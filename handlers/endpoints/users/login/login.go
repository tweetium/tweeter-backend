package login

import (
	"net/http"
	"time"

	"tweeter/db/models/user"
	handlerContext "tweeter/handlers/context"
	"tweeter/handlers/endpoints"
	"tweeter/handlers/endpoints/users"
	usersJWT "tweeter/handlers/endpoints/users/jwt"
	"tweeter/handlers/responses"
	"tweeter/handlers/util"
)

// Endpoint is the /api/users/login endpoint
var Endpoint = endpoints.Endpoint{
	Name:    "users#login",
	URL:     "/api/users/login",
	Handler: handleUserLogin,
	Methods: []string{http.MethodPost},
}

func handleUserLogin(req *http.Request, ctx handlerContext.Context) {
	type UserLoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginReq UserLoginReq
	ok := util.ParseBody(req, ctx, &loginReq)
	if !ok {
		return
	}

	loginUser, err := user.Find(user.FindByEmail{Email: loginReq.Email})
	if err != nil {
		switch err {
		case user.ErrInternalError:
			ctx.RenderErrorResponse(http.StatusInternalServerError, responses.ErrInternalError)
		case user.ErrUserNotFound:
			ctx.RenderErrorResponse(http.StatusBadRequest, users.ErrInvalidCredentials)
		default:
			// Logged as error because this indicates a programmer error, should fix the code if this happens
			ctx.Logger().WithError(err).Error("Uncaught error for user.Login")
			ctx.RenderErrorResponse(http.StatusInternalServerError, responses.ErrInternalError)
		}
		return
	}

	passwordMatches := loginUser.ComparePassword(loginReq.Password)
	if !passwordMatches {
		ctx.RenderErrorResponse(http.StatusBadRequest, users.ErrInvalidCredentials)
		return
	}

	// now the credentials are validated, grant identification token
	var tokenString string
	tokenString, err = usersJWT.GenerateToken(usersJWT.Claims{
		UserID: loginUser.ID,
	})
	if err != nil {
		ctx.RenderErrorResponse(http.StatusInternalServerError, responses.ErrInternalError)
		return
	}

	ctx.SetCookie(&http.Cookie{
		Name:    usersJWT.CookieName,
		Value:   tokenString,
		Expires: time.Now().AddDate(0, 1 /* month */, 0),
	})
	ctx.RenderResponse(http.StatusOK, responses.NewSuccessResponse(nil))
}
