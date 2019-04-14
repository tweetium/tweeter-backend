package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	handlerContext "tweeter/handlers/context"
	"tweeter/handlers/responses"
)

// ParseBody parses a request's body as the type provided, returning false if there was an error
// There is no need to render any response if there was an error, this function will render.
func ParseBody(req *http.Request, ctx handlerContext.Context, v interface{}) bool {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// This is unexpected (but possible), so let's log this internally here
		ctx.Logger().WithError(err).Warn("Failed to read request body")
		ctx.RenderErrorResponse(http.StatusBadRequest, responses.Error{
			Title: "Malformed Body", Detail: fmt.Sprintf("Failed to read request body"),
		})
		return false
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		ctx.RenderErrorResponse(http.StatusBadRequest, responses.ErrInvalidJSONBody)
		return false
	}

	return true
}
