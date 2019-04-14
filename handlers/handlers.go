package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	usersCreate "tweeter/handlers/endpoints/users/create"
)

// RunWebserver starts up the webserver and blocks until it is finished
func RunWebserver(port uint32) {
	r := mux.NewRouter()
	usersCreate.Endpoint.Attach(r)

	// Attach prometheus endpoint
	handler := promhttp.Handler()
	r.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	})

	logrus.WithField("port", port).Info("Http server started")
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		logrus.WithError(err).Fatal("Http server exited with error")
	}
}
