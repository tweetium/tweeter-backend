package endpoints

import (
	"net/http"

	"github.com/getsentry/raven-go"
	"github.com/gorilla/mux"

	"tweeter/handlers/context"
	"tweeter/handlers/middleware"
)

// HandlerFunc is type of function used to handle http requests and has
// additional endpoint.Context
type HandlerFunc func(*http.Request, context.Context)

// Endpoint is the type required to define an endpoint
type Endpoint struct {
	Name    string
	URL     string
	Handler HandlerFunc
	Methods []string
}

// Attach attaches the endpoint defined to the global http server
func (e Endpoint) Attach(r *mux.Router) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.New(e.Name, w, r)
		e.Handler(r, ctx)
	}
	handler = middleware.Log(handler)
	handler = middleware.Metrics(e.Name, handler)
	handler = raven.RecoveryHandler(handler)

	route := r.HandleFunc(e.URL, handler)
	if e.Methods != nil {
		route.Methods(e.Methods...)
	}
}
