package connection

import (
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

var c *gorm.DB

//GetConnection gives us a GORM DB object that we can use to talk to the database
func GetConnection() (*gorm.DB, error) {
	if c != nil {
		return c, nil
	}

	// conn is a GORM abstraction which internally wraps a connection pool
	conn, err := gorm.Open("postgres", "host=localhost port=6666 user=postgres dbname=postgres password=password sslmode=disable")
	if err != nil {
		return nil, err
	}

	// the size of this pool is set here. We cannot have more than 3 open connections to the database (see `pool`)
	conn.DB().SetMaxOpenConns(10)

	// optionall, we can opt-in for auto-preloading. That is, tell GORM to automatically load all related models
	// conn = conn.Set("gorm:auto_preload", true)

	// not thread safe, we don't care in this exercise
	c = conn

	return c, nil
}

//ClearDatabase removes the data from all the user tables in the database
func ClearDatabase() error {
	conn, err := GetConnection()
	if err != nil {
		return err
	}

	var tables []string
	err = conn.Table("pg_tables").Where("schemaname = 'public' and tablename != 'schema_migrations'").Pluck("tablename", &tables).Error
	if err != nil {
		return err
	}

	return conn.Exec("TRUNCATE TABLE " + strings.Join(tables, ",") + " CASCADE").Error
}

//RunMigrations applies the migrations found under /migrations
func RunMigrations() error {
	conn, err := GetConnection()
	if err != nil {
		return err
	}
	driver, err := postgres.WithInstance(conn.DB(), &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		panic(err)
	}
	return m.Up()
}
