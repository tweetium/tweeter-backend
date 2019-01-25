package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"tweeter/handlers/responses"
)

func render(w http.ResponseWriter, statusCode int, resp interface{}) {
	respBytes, err := json.Marshal(resp)
	if err != nil {
		captureError(err)
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(respBytes); err != nil {
		captureError(err)
	}
}

func renderErrors(w http.ResponseWriter, statusCode int, errors ...responses.Error) {
	render(w, statusCode, responses.NewErrorResponse(errors...))
}

func captureError(err error) {
	log.Printf("err() = %s", err)
}
