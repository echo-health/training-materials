package models

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

type PetKind string

func (k PetKind) String() string { return string(k) }

const (
	Cat PetKind = "cat"
	Dog PetKind = "dog"
)

type Person struct {
	Name    string
	Surname string
}

func (Person) TableName() string {
	return "persons"
}

type Pet struct {
	// you can optionally embed `gorm.Model` here, which will add some useful fields like ID, CreatedAt, etc.
	// the reason we don't do it is because we manage migrations separately and want to have more control
	// over the values of these fields
	Name      string
	Kind      string
	OwnerName sql.NullString
	Owner     *Person `gorm:"foreignkey:OwnerName"`
}

func (person *Person) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Surname", "Smith")
	return nil
}
