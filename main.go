package main

import (
	"fmt"
	"net/http"

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

	runWebserver(port)
}

func runWebserver(port uint32) {
	http.HandleFunc("/api/v1/users", logMiddleware(handlers.UsersHandler))

	logrus.WithFields(logrus.Fields{"port": port}).Info("Http server started")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Http server exited with error")
	}
}

func logMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL,
		}).Debug("Received http request")
		handler(w, r)
	}
}
