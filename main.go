package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"tweeter/util"

	_ "github.com/lib/pq" // required driver
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file" // required driver
)

func main() {
	port := util.MustGetEnvUInt32("PORT")
	dbURL := util.MustGetEnv("DATABASE_URL")

	err := migrateDatabase(dbURL)
	if err != nil {
		log.Fatalf("Failed to migrate database: %s", err)
	}

	runWebserver(port)
}

func runWebserver(port uint32) {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "Hello, world!\n")
	}
	http.HandleFunc("/hello", helloHandler)

	log.Printf("Listening at :%d!..", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func migrateDatabase(path string) error {
	db, err := sql.Open("postgres", path)
	defer db.Close()
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
	log.Printf("Currently at version: %d, dirty: %t, err: %s", version, dirty, err)
	m.Close()
	log.Printf("Done running migrations!")

	return nil
}
