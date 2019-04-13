package main

import (
	_ "github.com/lib/pq"                     // required driver
	_ "github.com/mattes/migrate/source/file" // required driver
	"github.com/sirupsen/logrus"

	"tweeter/db"
	"tweeter/handlers"
	"tweeter/log"
	"tweeter/util"
)

func main() {
	port := util.MustGetEnvUInt32("PORT")
	dbURL := util.MustGetEnv("DATABASE_URL")

	log.Init()

	err := db.Init(dbURL)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to initialize database")
	}

	handlers.RunWebserver(port)
}
