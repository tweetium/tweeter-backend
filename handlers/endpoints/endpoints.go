package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"

	"tweeter/handlers/middleware"
)

// Endpoint is the type required to define an endpoint
type Endpoint struct {
	URL     string
	Handler http.HandlerFunc
	Methods []string
}

// Attach attaches the endpoint defined to the global http server
func (e Endpoint) Attach(r *mux.Router) {
	route := r.HandleFunc(e.URL, middleware.Log(e.Handler))
	if e.Methods != nil {
		route.Methods(e.Methods...)
	}
}
