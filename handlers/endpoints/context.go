package endpoints

import (
	"net/http"

	"tweeter/handlers/render"
	"tweeter/handlers/responses"
)

// Context is a struct with additional utility for the handler func
type Context struct {
	endpoint *Endpoint

	responseWriter http.ResponseWriter
	request        *http.Request
}

func (c Context) RenderResponse(statusCode int, resp interface{}) {
	render.Response(c.endpoint.Name, c.responseWriter, statusCode, resp)
}

func (c Context) RenderErrorResponse(statusCode int, errors ...responses.Error) {
	render.ErrorResponse(c.endpoint.Name, c.responseWriter, statusCode, errors...)
}
