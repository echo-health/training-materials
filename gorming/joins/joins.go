package main

import (
	"database/sql"
	"fmt"

	"github.com/echo-health/training-materials/gorming/connection"
	"github.com/echo-health/training-materials/gorming/models"
)

type results []struct {
	OwnerName string
	PetName   string
	PetKind   string
}

func printRows(rows results) {
	for _, row := range rows {
		if row.OwnerName == "" {
			fmt.Printf("%s (the %s) has no owner!\n", row.PetName, row.PetKind)
		} else if row.PetName != "" {
			fmt.Printf("%s has %s, which is a %s.\n", row.OwnerName, row.PetName, row.PetKind)
		} else {
			fmt.Printf("%s has no pets.\n", row.OwnerName)
		}
	}
}

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
	peter := &models.Person{Name: "Peter"}
	for _, person := range []*models.Person{lucas, jesus, peter} {
		if err = conn.Create(person).Error; err != nil {
			panic(err)
		}
	}

	// create two pets
	cat := &models.Pet{
		Name:      "Luna",
		Kind:      models.Cat.String(),
		OwnerName: sql.NullString{String: lucas.Name, Valid: true},
	}
	dog := &models.Pet{
		Name:      "Bruce",
		Kind:      models.Dog.String(),
		OwnerName: sql.NullString{String: lucas.Name, Valid: true},
	}
	bowie := &models.Pet{
		Name: "Bowie",
		Kind: models.Dog.String(),
	}
	for _, pet := range []*models.Pet{cat, dog, bowie} {
		if err = conn.Create(pet).Error; err != nil {
			panic(err)
		}
	}

	var rows results

	fmt.Println("-----------------LEFT JOIN-----------------")

	// LEFT JOIN: PERSONS -> PETS
	if err = conn.Model(&models.Person{}).Joins("LEFT JOIN pets ON pets.owner_name = persons.name").Select("persons.name as owner_name, pets.name as pet_name, pets.kind as pet_kind").Scan(&rows).Error; err != nil {
		panic(err)
	}
	printRows(rows)

	// RIGHT JOIN: PERSONS <- PETS
	fmt.Println("-----------------RIGHT JOIN----------------")

	if err = conn.Model(&models.Person{}).Joins("RIGHT JOIN pets ON pets.owner_name = persons.name").Select("persons.name as owner_name, pets.name as pet_name, pets.kind as pet_kind").Scan(&rows).Error; err != nil {
		panic(err)
	}
	printRows(rows)

	// INNER JOIN: PERSONS <-> PETS
	fmt.Println("-----------------INNER JOIN----------------")

	if err = conn.Model(&models.Person{}).Joins("INNER JOIN pets ON pets.owner_name = persons.name").Select("persons.name as owner_name, pets.name as pet_name, pets.kind as pet_kind").Scan(&rows).Error; err != nil {
		panic(err)
	}
	printRows(rows)

	// FULL OUTER JOIN: PERSONS AND PETS
	fmt.Println("---------------FULL OUTER JOIN-------------")

	if err = conn.Model(&models.Person{}).Joins("FULL OUTER JOIN pets ON pets.owner_name = persons.name").Select("persons.name as owner_name, pets.name as pet_name, pets.kind as pet_kind").Scan(&rows).Error; err != nil {
		panic(err)
	}
	printRows(rows)

	// INNER JOIN: PERSONS <-> PETS (only cats)
	fmt.Println("-------------FILTERED INNER JOIN-----------")

	// you can always filter on any of the tables
	if err = conn.Model(&models.Person{}).Joins("INNER JOIN pets ON pets.owner_name = persons.name").Select("persons.name as owner_name, pets.name as pet_name, pets.kind as pet_kind").Where("pets.kind = ?", models.Cat.String()).Scan(&rows).Error; err != nil {
		panic(err)
	}
	printRows(rows)
}
