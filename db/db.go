package db

import (
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
)

// DB is the main database
var DB *sqlx.DB

// Init initializes and migrates the database given the provided connect url
func Init(dbURL string) (err error) {
	DB, err = sqlx.Connect("postgres", dbURL)
	if err != nil {
		return
	}

	err = migrateDatabase(dbURL)
	log.Printf("Connected to database at: %s", dbURL)
	return
}

func migrateDatabase(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	log.Printf("Running migrations..")
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations/",
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	m.Up()
	version, dirty, err := m.Version()
	if err != nil {
		log.Printf("Err when running migration.Version(): %s", err)
	} else {
		log.Printf("Currently at version: %d, dirty: %t", version, dirty)
	}
	m.Close()
	log.Printf("Done running migrations!")

	return nil
}
