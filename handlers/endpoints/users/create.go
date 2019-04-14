package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"tweeter/db/models/user"
	"tweeter/handlers/context"
	"tweeter/handlers/endpoints"
	"tweeter/handlers/responses"
)

// CreateEndpoint is the /api/v1/users/ create endpoint
var CreateEndpoint = endpoints.Endpoint{
	Name:    "users#create",
	URL:     "/api/v1/users",
	Handler: handleUserCreate,
	Methods: []string{http.MethodPost},
}

func handleUserCreate(req *http.Request, ctx context.Context) {
	if req.Method != http.MethodPost {
		// This shouldn't be possible given that the route only accepts POST requests
		ctx.Logger().WithField("method", req.Method).Error("Invalid method for users#create")
		ctx.RenderErrorResponse(http.StatusInternalServerError, responses.ErrInternalError)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// This is unexpected (but possible), so let's log this internally here
		ctx.Logger().WithError(err).Warn("Failed to read request body")
		ctx.RenderErrorResponse(http.StatusBadRequest, responses.Error{
			Title: "Malformed Body", Detail: fmt.Sprintf("Failed to read request body"),
		})
		return
	}

	type UserCreateReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var createReq UserCreateReq
	err = json.Unmarshal(body, &createReq)
	if err != nil {
		ctx.RenderErrorResponse(http.StatusBadRequest, ErrInvalidBody)
		return
	}

	var newUser user.User
	newUser, err = user.Create(createReq.Email, createReq.Password)
	if err != nil {
		switch err {
		case user.ErrInternalError:
			ctx.RenderErrorResponse(http.StatusInternalServerError, responses.ErrInternalError)
		case user.ErrPasswordTooShort:
			ctx.RenderErrorResponse(http.StatusBadRequest, ErrPasswordTooShort)
		case user.ErrUserEmailAlreadyExists:
			ctx.RenderErrorResponse(http.StatusBadRequest, ErrEmailAlreadyExists(createReq.Email))
		default:
			// Logged as error because this indicates a programmer error, should fix the code if this happens
			ctx.Logger().WithError(err).Error("Uncaught error for user.Create")
			ctx.RenderErrorResponse(http.StatusInternalServerError, responses.ErrInternalError)
		}
		return
	}

	ctx.RenderResponse(http.StatusOK, responses.NewSuccessResponse(struct {
		ID user.ID
	}{
		newUser.ID,
	}))
}
