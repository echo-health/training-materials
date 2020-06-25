package main

import (
	"fmt"

	"github.com/echo-health/training-materials/gorming/connection"
	"github.com/echo-health/training-materials/gorming/models"
)

func main() {
	conn, err := connection.GetConnection()
	if err != nil {
		panic(err)
	}

	// apply migrations
	connection.RunMigrations()

	// clean up tables when we're done
	defer connection.ClearDatabase()

	// create one person
	lucas := &models.Person{Name: "Lucas"}
	jesus := &models.Person{Name: "Jesus"}
	for _, person := range []*models.Person{lucas, jesus} {
		if err = conn.Create(person).Error; err != nil {
			panic(err)
		}
	}

	// create two pets
	cat := &models.Pet{
		Name:      "Luna",
		Kind:      models.Cat.String(),
		OwnerName: lucas.Name,
	}
	dog := &models.Pet{
		Name:      "Bruce",
		Kind:      models.Dog.String(),
		OwnerName: lucas.Name,
	}
	for _, pet := range []*models.Pet{cat, dog} {
		if err = conn.Create(pet).Error; err != nil {
			panic(err)
		}
	}

	// get all persons
	var persons []*models.Person
	if err = conn.Find(&persons).Error; err != nil {
		panic(err)
	}
	fmt.Printf("Found %d persons (%s)\n", len(persons), persons)

	// both have a surname, but we did not provide it! This has been done by a GORM hook (https://gorm.io/docs/hooks.html)
	// there are multiple hooks you can use. It's a good opportunity to set some default fields for your models
	fmt.Printf("Full names: %s %s and %s %s\n", persons[0].Name, persons[0].Surname, persons[1].Name, persons[1].Surname)

	// get cat
	var luna models.Pet
	if err = conn.Where("kind = ?", models.Cat).First(&luna).Error; err != nil {
		panic(err)
	}
	fmt.Printf("Found %s and it is a %s\n", luna.Name, luna.Kind)

	// check if Luna has an owner
	fmt.Printf("Does Luna have an owner? %t\n", luna.Owner != nil)

	// but it should have one! what's up? Relationships need to be either explicitly loaded or tell GORM to always do it
	// if we explicitly ask for it
	if err = conn.Where("kind = ?", models.Cat).Preload("Owner").First(&luna).Error; err != nil {
		panic(err)
	}

	// do the check again
	fmt.Printf("Does Luna have an owner? %t\n", luna.Owner != nil)

	// the alternative is to always tell gorm to preload relationships. Let's try with Bruce
	conn = conn.Set("gorm:auto_preload", true)

	// note that below we don't have any explicit
	var bruce models.Pet
	if err = conn.Where("kind = ?", models.Cat).First(&bruce).Error; err != nil {
		panic(err)
	}

	// ding ding
	fmt.Printf("Does Bruce have an owner? %t\n", bruce.Owner != nil)

	// the good thing about the preload option is that you can either do it on a per-query basis:
	// 	eg: conn.Set("gorm:auto_preload", true).Find(&pets)
	// or you can set it in the function you use to get your connection object, meaning everything will be auto-preloaded (see `GetConnection`)

	// if we just wanted the pet names
	var petNames []string
	if err = conn.Model(&models.Pet{}).Pluck("name", &petNames).Error; err != nil {
		panic(err)
	}

	// if we need a subset of the columns (not all of them)
	var columns []struct {
		Name string
		Kind string
	}
	if err = conn.Model(&models.Pet{}).Select("name, kind").Scan(&columns).Error; err != nil {
		panic(err)
	}

	// finally, let's see how pagination works. You should use this anywhere you're implementing a `List` method
	// and exposing it to other clients who don't necessarily need to know about how many rows to expect
	nCats := 100
	for nCats > 0 {
		cat := &models.Pet{Name: fmt.Sprintf("garfield-%d", nCats), Kind: models.Cat.String(), OwnerName: jesus.Name}
		if err = conn.Create(cat).Error; err != nil {
			panic(err)
		}
		nCats--
	}

	// now let's paginate over the results. In our services, you might see that we pass a `Token` (pagination token), which
	// is a string and then gets converted into an integer. The reason it is a string rather than just an integer is to
	// allow for flexibility if the pagination system changes, which could happen if we switched postgres with a different
	// database or if we re-defined how the data is sorted (eg: we could use the pet Name as a token, if it was unique)
	offset := 0
	for {
		var partialResults []*models.Pet
		if err = conn.Where("name ~ ?", "garfield").Offset(offset).Limit(33).Find(&partialResults).Error; err != nil {
			panic(err)
		}
		offset += len(partialResults)

		fmt.Printf("Got %d results, offset is now %d\n", len(partialResults), offset)

		// as soon as we get less results than we asked for, we know we're done
		if len(partialResults) < 33 {
			break
		}
	}
}
