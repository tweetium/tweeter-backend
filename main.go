package main

import (
	_ "github.com/lib/pq"                     // required driver
	_ "github.com/mattes/migrate/source/file" // required driver
	"github.com/sirupsen/logrus"

	"tweeter/db"
	"tweeter/handlers"
	"tweeter/util"
)

func main() {
	port := util.MustGetEnvUInt32("PORT")
	dbURL := util.MustGetEnv("DATABASE_URL")

	err := db.Init(dbURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Failed to initialize database")
	}

	handlers.RunWebserver(port)
}
