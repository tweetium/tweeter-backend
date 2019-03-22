package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"

	"tweeter/handlers/middleware"
	"tweeter/handlers/render"
	"tweeter/handlers/responses"
)

// Context is a struct with additional utility for the handler func
type Context struct {
	// RenderResponse renders a response from the resp provided
	RenderResponse func(http.ResponseWriter, int, interface{})
	// RenderErrorResponse renders the error response with the status code provided
	RenderErrorResponse func(http.ResponseWriter, int, ...responses.Error)
}

// HandlerFunc is type of function used to handle http requests and has
// additional endpoint.Context
type HandlerFunc func(http.ResponseWriter, *http.Request, Context)

// Endpoint is the type required to define an endpoint
type Endpoint struct {
	Name    string
	URL     string
	Handler HandlerFunc
	Methods []string
}

// Attach attaches the endpoint defined to the global http server
func (e Endpoint) Attach(r *mux.Router) {
	ctx := Context{
		RenderResponse: func(w http.ResponseWriter, statusCode int, resp interface{}) {
			render.Response(e.Name, w, statusCode, resp)
		},
		RenderErrorResponse: func(w http.ResponseWriter, statusCode int, errors ...responses.Error) {
			render.ErrorResponse(e.Name, w, statusCode, errors...)
		},
	}

	handler := func(w http.ResponseWriter, r *http.Request) { e.Handler(w, r, ctx) }
	handler = middleware.Log(handler)
	handler = middleware.Metrics(e.Name, handler)

	route := r.HandleFunc(e.URL, handler)
	if e.Methods != nil {
		route.Methods(e.Methods...)
	}
}
