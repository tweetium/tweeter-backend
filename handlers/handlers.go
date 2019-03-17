package handlers

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"tweeter/handlers/endpoints/users"
	"tweeter/handlers/middleware"
)

// RunWebserver starts up the webserver and blocks until it is finished
func RunWebserver(port uint32) {
	http.HandleFunc("/api/v1/users", middleware.Log(users.Handler))

	logrus.WithFields(logrus.Fields{"port": port}).Info("Http server started")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Http server exited with error")
	}
}
