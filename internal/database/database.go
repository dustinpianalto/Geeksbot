package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/dustinpianalto/geeksbot/internal/database/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"
)

var (
	db *sql.DB
)

func ConnectDatabase(dbConnString string) {
	s := bindata.Resource(migrations.AssetNames(),
		func(name string) ([]byte, error) {
			return migrations.Asset(name)
		})
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatal(fmt.Errorf("Can't connect to the database: %w", err))
	}
	log.Println("Database Connected.")
	db.SetMaxOpenConns(75)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(300)
	d, err := bindata.WithInstance(s)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot load migrations: %w", err))
	}
	instance, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(fmt.Errorf("cannot create db driver: %w", err))
	}
	m, err := migrate.NewWithInstance("go-bindata", d, "postgres", instance)
	if err != nil {
		log.Fatal(fmt.Errorf("cannot create migration instance: %w", err))
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrations Run")
}
