package create

import (
	"net/http"

	"tweeter/db/models/user"
	handlerContext "tweeter/handlers/context"
	"tweeter/handlers/endpoints"
	"tweeter/handlers/endpoints/users"
	"tweeter/handlers/responses"
	"tweeter/handlers/util"
)

// Endpoint is the /api/users/ create endpoint
var Endpoint = endpoints.Endpoint{
	Name:    "users#create",
	URL:     "/api/users",
	Handler: handleUserCreate,
	Methods: []string{http.MethodPost},
}

func handleUserCreate(req *http.Request, ctx handlerContext.Context) {
	type UserCreateReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var createReq UserCreateReq
	ok := util.ParseBody(req, ctx, &createReq)
	if !ok {
		return
	}

	newUser, err := user.Create(createReq.Email, createReq.Password)
	if err != nil {
		switch err {
		case user.ErrInternalError:
			ctx.RenderErrorResponse(http.StatusInternalServerError, responses.ErrInternalError)
		case user.ErrPasswordTooShort:
			ctx.RenderErrorResponse(http.StatusBadRequest, users.ErrPasswordTooShort)
		case user.ErrUserEmailAlreadyExists:
			ctx.RenderErrorResponse(http.StatusBadRequest, users.ErrEmailAlreadyExists(createReq.Email))
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
