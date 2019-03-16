package handlers

import (
	"encoding/json"
	"net/http"
	"tweeter/handlers/responses"

	"github.com/sirupsen/logrus"
)

func render(w http.ResponseWriter, statusCode int, resp interface{}) {
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
}

func renderErrors(w http.ResponseWriter, statusCode int, errors ...responses.Error) {
	render(w, statusCode, responses.NewErrorResponse(errors...))
}
