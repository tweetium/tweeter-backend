package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"tweeter/handlers/endpoints/users"
)

// RunWebserver starts up the webserver and blocks until it is finished
func RunWebserver(port uint32) {
	r := mux.NewRouter()
	users.CreateEndpoint.Attach(r)

	logrus.WithFields(logrus.Fields{"port": port}).Info("Http server started")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Fatal("Http server exited with error")
	}
}
