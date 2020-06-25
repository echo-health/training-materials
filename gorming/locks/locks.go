package main

import (
	"fmt"
	"time"

	"math/rand"

	"github.com/echo-health/training-materials/gorming/connection"
	"github.com/echo-health/training-materials/gorming/models"
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
	connection.ClearDatabase()
	defer connection.ClearDatabase()

	// create one person
	jesus := &models.Person{Name: "Jesus"}
	for _, person := range []*models.Person{jesus} {
		if err = conn.Create(person).Error; err != nil {
			panic(err)
		}
	}

	// this can be a very dense and complex topic to cover in detail (we can leave that for another training course). Let's
	// simply focus on what we do in ECHO

	// when working in any 'multi-threaded' environment (a web server is a good example, unless you decide to run it with a single
	// thread, good luck), there are occassions when you want to make sure you have **EXCLUSIVE** access to the data you are
	// modifying. The usual requirements are:
	//	- you are ok with others reading the data you're modifying, while you modify it
	//  - you don't want anyone else to modify the same data, while you're modifying it
	//  - anyone wanting to modify the same data, should always get the **latest** version of the data when it's available to be modified

	// databases provide different levels of locking in which those requirements above might be eased or tightened. In our case, whenever
	// we want to do the above, we use postgres' `SELECT FOR UPDATE` which essentially tells postgres to lock the rows we're selecting for us.
	// Locks are only held for the duration of a transaction. This means:
	//  - be careful with leaking transactions, you could end up with your database in an unusuable state (not really true, there are workarounds)
	//  - make sure you always open a transaction if you want to hold a lock for some time
	folk1, folk2 := make(chan bool), make(chan bool)
	func() {
		tx := conn.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		var person models.Person
		if err = tx.Set("gorm:query_option", "FOR UPDATE").Where("name = ?", jesus.Name).First(&person).Error; err != nil {
			panic(err)
		}

		// now person is locked for just for me! While this transaction is open, the lock will be held, so I know no one is
		// going to overwrite what I do here

		// let's simulate the case where someone else comes in and wants to change my name

		go func() {
			anotherConn, err := connection.GetConnection()
			if err != nil {
				panic(err)
			}

			fmt.Println("Someone is waiting to update the name")
			// this will wait until the lock `tx` acquired is released
			if err = anotherConn.Model(&models.Person{}).Where("name = ?", jesus.Name).Update("name", "Sammy the Snail").Error; err != nil {
				panic(err)
			}
			folk1 <- true
		}()

		// yet another folk comes in, and wants to change my name..but it might not succeed
		go func() {
			anotherConn, err := connection.GetConnection()
			if err != nil {
				panic(err)
			}

			fmt.Println("Someone else is also waiting to update the surname")
			// this will wait until the lock `tx` acquired is released
			if err = anotherConn.Model(&models.Person{}).Where("name = ?", jesus.Name).Update("surname", "Martin").Error; err != nil {
				panic(err)
			}
			folk2 <- true
		}()

		// and finally, someone just wants to check the current surname
		go func() {
			anotherConn, err := connection.GetConnection()
			if err != nil {
				panic(err)
			}

			fmt.Println("And another one just wants to know the surname")
			// this one just wants to read, and that's allowed
			var surnames []string
			if err = anotherConn.Model(&models.Person{}).Where("name = ?", jesus.Name).Pluck("surname", &surnames).Error; err != nil {
				panic(err)
			}
			fmt.Printf("The surname is %s\n", surnames[0])
		}()

		// let's simulate some work happening
		time.Sleep(2 * time.Second)

		// then we update the surname
		if err = tx.Model(&models.Person{}).Where("name = ?", jesus.Name).Update("surname", "Hernandez").Error; err != nil {
			panic(err)
		}

		// at this point the lock will be released
		tx.Commit()
	}()

	// wait until these two have finished
	<-folk1
	<-folk2

	// now, think about the difference possibilities here. Both folks were competing between each other to acquire the lock. Depending
	// on who acquired it first, we'll see a different result here!
	if err = conn.First(&jesus).Error; err != nil {
		panic(err)
	}
	fmt.Println("------------------------------")
	fmt.Printf("The result is: name=%s, surname=%s\n", jesus.Name, jesus.Surname)
}
