package context

import (
	"encoding/json"
	"net/http"
	"strconv"

	"tweeter/handlers/responses"
	"tweeter/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// Context is a struct with additional utility for the handler func
type Context struct {
	endpointName string

	responseWriter http.ResponseWriter
	request        *http.Request
}

// New creates a Context
func New(endpointName string, responseWriter http.ResponseWriter, request *http.Request) Context {
	return Context{
		endpointName,
		responseWriter,
		request,
	}
}

// Logger returns a pre-populated logrus.Entry to log with
func (c Context) Logger() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"endpointName": c.endpointName,
		"request":      c.request,
	})
}

// SetCookie sets a cookie for the current request
func (c Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.responseWriter, cookie)
}

// RenderResponse renders a response from type provided
func (c Context) RenderResponse(statusCode int, resp interface{}) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		// This is an error because resp is controlled by the programmer and
		// should be correct in all situations
		c.Logger().WithError(err).Error("Resp passed to render was not json.Marshalable")
		return
	}

	c.responseWriter.WriteHeader(statusCode)
	if _, err := c.responseWriter.Write(respBytes); err != nil {
		c.Logger().WithError(err).Warn("Failed to write bytes to http.ResponseWriter")
	}

	metrics.APIResponses.
		With(prometheus.Labels{
			"endpointName": c.endpointName,
			"statusCode":   strconv.FormatInt(int64(statusCode), 10),
		}).Inc()
}

// RenderErrorResponse renders a standard error with the errors provided
func (c Context) RenderErrorResponse(statusCode int, errors ...responses.Error) {
	for _, error := range errors {
		metrics.APIResponseErrors.
			With(prometheus.Labels{
				"endpointName": c.endpointName,
				"errorTitle":   error.Title,
			}).Inc()
	}

	c.RenderResponse(statusCode, responses.NewErrorResponse(errors...))
}
