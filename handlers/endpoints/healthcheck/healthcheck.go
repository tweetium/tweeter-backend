package healthcheck

import (
	"net/http"

	handlerContext "tweeter/handlers/context"
	"tweeter/handlers/endpoints"
	"tweeter/handlers/responses"
)

// Endpoint is the /healthcheck endpoint
var Endpoint = endpoints.Endpoint{
	Name:    "healthcheck",
	URL:     "/healthcheck",
	Handler: handleHealthcheck,
	Methods: []string{http.MethodGet},
}

func handleHealthcheck(_req *http.Request, ctx handlerContext.Context) {
	ctx.RenderResponse(http.StatusOK, responses.NewSuccessResponse(nil))
}
