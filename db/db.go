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

// DB is a postgres Database connection / connection pool
var DB Database

var db *sqlx.DB
var testTx *sqlx.Tx

// Database is an abstraction over a sqlx.DB / sqlx.Tx (transaction)
// a sqlx.DB is a connection pool and a Tx is a connection that has a transaction running on it
type Database interface {
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	MustExec(query string, args ...interface{}) sql.Result
}

// Init initializes and migrates the database given the provided connect url
func Init(dbURL string) (err error) {
	db, err = sqlx.Connect("postgres", dbURL)
	DB = db
	if err != nil {
		return
	}

	err = migrateDatabase(dbURL)
	logrus.WithFields(logrus.Fields{"url": dbURL}).Info("Connected to database")
	return
}

func migrateDatabase(dbURL string) error {
	// Use a different connection for migrations as we close the source below
	// An attempt was made to re-use the existing db opened in Init, but it failed
	// due to the migration.Up() failing due to "can't acquire lock"
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

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Error running migration.Up()")
	}

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

// BeginTransactionForTests creates a new transaction for the test
// and injects it as the DB connection
func BeginTransactionForTests() {
	var err error
	testTx, err = db.Beginx()
	DB = testTx
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Failed to begin DB transaction for tests")
	}
}

// RollbackTransactionForTests cleans up the database for tests
func RollbackTransactionForTests() {
	if testTx == nil {
		logrus.Fatal("Cannot TeardownForTests without testTx (did you call InitForTests?)")
	}

	err := testTx.Rollback()
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Failed to rollback DB transaction for tests")
	}

	DB = db
}

var migrationPath = "file://migrations/"
