package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"

	"tweeter/metrics"
)

// Metrics handles generic metrics around a request that a Handler gets
func Metrics(endpointName string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metrics.APIReceived.With(prometheus.Labels{"endpointName": endpointName}).Inc()
		handler(w, r)
	}
}
