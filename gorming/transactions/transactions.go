package main

import (
	"database/sql"
	"fmt"
	"time"

	"math/rand"

	"github.com/echo-health/training-materials/gorming/connection"
	"github.com/echo-health/training-materials/gorming/models"
	"github.com/jinzhu/gorm"
)

func main() {
	rand.Seed(time.Now().Unix())

	conn, err := connection.GetConnection()
	if err != nil {
		panic(err)
	}

	// apply migrations
	connection.RunMigrations()

	// clean up tables when we're done
	defer connection.ClearDatabase()

	// if you just need to create a new object (add a single row), you can just use the connection object
	lucas := &models.Person{Name: "Lucas"}
	if err = conn.Create(lucas).Error; err != nil {
		panic(err)
	}

	// you can read it immediately after
	var person models.Person
	if err = conn.First(&person).Error; err != nil {
		panic(err)
	}
	fmt.Printf("Found one person! (%s)\n", person.Name)

	// now imagine you're runing some piece of code which needs to create a person, then do some important stuff
	// and finally, if that stuff runs successfully, create a pet. Someone that does not know about transactions
	// might be tempted to do this (wrapping it in a function so we can clean up easily!)
	func() {
		defer connection.ClearDatabase()
		jesus := &models.Person{Name: "Jesus"}
		if err = conn.Create(jesus).Error; err != nil {
			panic(err)
		}

		// the important stuff
		isNumberEven := rand.Int()%2 == 0

		if isNumberEven {
			// if we're successful, create a cat
			cat := &models.Pet{
				Name:      "Luna",
				Kind:      models.Cat.String(),
				OwnerName: sql.NullString{String: jesus.Name, Valid: true},
			}
			if err = conn.Create(cat).Error; err != nil {
				panic(err)
			}
		} else {
			// if we're not, we want to clean up. Delete the person
			if err = conn.Delete(jesus).Error; err != nil {
				panic(err)
			}
		}
	}()

	// obviously the above is not ideal. Many things can go wrong after we create the person:
	// - the cat might fail to be created
	// - isNumberEven could panic (in this case, it won't, but think of any other func call)
	// - we might fail to delete the person if things go wrong

	// how do we solve it? Well, that's when database transactions come handy! They allow us to do
	// several operations atomically: either all succeed or all fail. GORM has a simple API to work
	// with transactions
	func() {
		defer connection.ClearDatabase()

		// we encapsulate the transaction block (in case the `ClearDatabase` call panics)
		func() {

			// open a transaction (tx is just another connection object)
			tx := conn.Begin()

			// make sure that if anything panics, we'll recover
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()

			// from now on, **all** database operations that you want to happen atomically, need to be
			// done using `tx` (if you use `conn`, you'll do them outside of the transaction!)
			jesus := &models.Person{Name: "Jesus"}
			if err = tx.Create(jesus).Error; err != nil {
				panic(err)
			}

			// remember that within a transaction, you can read all your uncomminted changes! (make sure you use `tx`)
			var person models.Person
			if err = tx.First(&person).Error; err != nil {
				panic(err)
			}
			fmt.Printf("Hello, I'm %s\n", person.Name)

			// BUT the key bit (and this is made sure by the database, thanks to the transaction isolation level), is that
			// other transactions won't be able to see the changes made here until they are commited!
			done := make(chan bool)
			go func() {
				otherConn, err := connection.GetConnection()
				if err != nil {
					panic(err)
				}
				var person models.Person
				err = otherConn.First(&person).Error
				if gorm.IsRecordNotFoundError(err) {
					fmt.Println("Oops! there are no persons! What?")
				}
				done <- true
			}()

			// do not continue after the other transaction checks if there are any persons
			<-done

			// the important stuff
			isNumberEven := rand.Int()%2 == 0

			if isNumberEven {
				// if we're successful, create a cat
				cat := &models.Pet{
					Name:      "Luna",
					Kind:      models.Cat.String(),
					OwnerName: sql.NullString{String: jesus.Name, Valid: true},
				}
				if err = tx.Create(cat).Error; err != nil {
					// this will be captured by the `defer`
					panic(err)
				}
				tx.Commit()
			} else {
				tx.Rollback()
			}
		}()
		// at this point the transaction has already commited or rolled back (based on `isNumberEven`)
		// lets see what we end up with (now the changes are visible to everyone, see how we use `conn` again)

		var person models.Person
		if err = conn.First(&person).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				fmt.Println("No person found. Which means there must be also no pet (as they are either all created or not)")

				var pet models.Pet
				err = conn.First(&pet).Error
				if gorm.IsRecordNotFoundError(err) {
					fmt.Println("Cool, as expected, there are no pets either")
				}
			}
		} else {
			fmt.Printf("We found %s\n", person.Name)

			var pet models.Pet
			if err = conn.First(&pet).Error; err != nil {
				panic(err)
			}
			fmt.Printf("And a pet: %s\n", pet.Name)
		}

	}()

}
