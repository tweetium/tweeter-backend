package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// APIReceived records the number of requests received to an API endpoint.
	APIReceived = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_received",
			Help: "Number of requests received to an API endpoint.",
		},
		[]string{"endpointName"},
	)

	// APIResponses records the number of responses returned from an API endpoint.
	APIResponses = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_responses",
			Help: "Number of responses from an API endpoint.",
		},
		[]string{"endpointName", "statusCode"},
	)

	// APIResponseErrors records the number of errors returned from an API endpoint.
	APIResponseErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_response_errors",
			Help: "Number of response.Errors by title.",
		},
		[]string{"endpointName", "errorTitle"},
	)
)

func init() {
	prometheus.MustRegister(APIReceived)
	prometheus.MustRegister(APIResponses)
	prometheus.MustRegister(APIResponseErrors)
}
