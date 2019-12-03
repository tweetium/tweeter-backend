package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"tweeter/handlers/endpoints/healthcheck"
	usersCreate "tweeter/handlers/endpoints/users/create"
	usersLogin "tweeter/handlers/endpoints/users/login"
)

// RunWebserver starts up the webserver and blocks until it is finished
func RunWebserver(port uint32) {
	router := GetRouter()

	logrus.WithField("port", port).Info("Http server started")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	if err != nil {
		logrus.WithError(err).Fatal("Http server exited with error")
	}
}

// GetRouter is a helper function to get the router with all the endpoints configured.
// This is used for testing that the server is setup properly.
func GetRouter() *mux.Router {
	r := mux.NewRouter()
	usersCreate.Endpoint.Attach(r)
	usersLogin.Endpoint.Attach(r)
	healthcheck.Endpoint.Attach(r)

	// Attach prometheus endpoint
	handler := promhttp.Handler()
	r.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	return r
}
