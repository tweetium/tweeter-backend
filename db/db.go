package db

import (
	"database/sql"
	"tweeter/util"

	_ "github.com/lib/pq"                     // required driver
	_ "github.com/mattes/migrate/source/file" // required driver
	"github.com/sirupsen/logrus"

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
	logrus.WithFields(logrus.Fields{"url": dbURL}).Info("Connected to database")
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

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", driver)
	if err != nil {
		return err
	}

	m.Up()
	version, dirty, err := m.Version()
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Error running migration.Version()")
	}
	m.Close()
	logrus.WithFields(logrus.Fields{
		"version": version,
		"dirty":   dirty,
	}).Info("Finished running migrations!")

	return nil
}

// InitForTests initialize the database for tests, requires different migration path
// since tests are run in the working directory they're located in.
func InitForTests() {
	migrationPath = "file:///app/migrations/"
	err := Init(util.MustGetEnv("DATABASE_URL"))
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Failed to initialize DB for tests")
	}
}

var migrationPath = "file://migrations/"
