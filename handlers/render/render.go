package render

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tweeter/handlers/responses"
	"tweeter/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Response renders a response from the resp provided
func Response(endpointName string, w http.ResponseWriter, statusCode int, resp interface{}) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		// This is an error because resp is controlled by the programmer and
		// should be correct in all situations
		logrus.WithFields(logrus.Fields{"err": err}).Error("Resp passed to render was not json.Marshalable")
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(respBytes); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Warn("Failed to write bytes to http.ResponseWriter")
	}

	metrics.APIResponses.
		With(prometheus.Labels{
			"endpointName": endpointName,
			"statusCode":   strconv.FormatInt(int64(statusCode), 10),
		}).Inc()
}

// ErrorResponse renders the error response with the status code provided
func ErrorResponse(endpointName string, w http.ResponseWriter, statusCode int, errors ...responses.Error) {
	for _, error := range errors {
		metrics.APIResponseErrors.
			With(prometheus.Labels{
				"endpointName": endpointName,
				"errorTitle":   error.Title,
			}).Inc()
	}

	Response(endpointName, w, statusCode, responses.NewErrorResponse(errors...))
}
