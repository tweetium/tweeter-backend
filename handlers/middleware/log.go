package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Log handles generic logging around a request that a Handler gets
func Log(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL,
		}).Debug("Received http request")
		handler(w, r)
	}
}
