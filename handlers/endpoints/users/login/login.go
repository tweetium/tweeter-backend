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

	"github.com/sirupsen/logrus"
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
			ctx.RenderInternalErrorResponse(err, logrus.WarnLevel, "user.Find encountered internal error")
		case user.ErrUserNotFound:
			ctx.RenderErrorResponse(http.StatusBadRequest, users.ErrInvalidCredentials)
		default:
			// Logged as error because this indicates a programmer error, should fix the code if this happens
			ctx.RenderInternalErrorResponse(err, logrus.ErrorLevel, "Uncaught error for user.Login")
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
		// This indicates a problem with the jwt secrets initialization - should be fixed by dev
		ctx.RenderInternalErrorResponse(err, logrus.ErrorLevel, "usersJWT.GenerateToken failed")
		return
	}

	ctx.SetCookie(&http.Cookie{
		Name:    usersJWT.CookieName,
		Value:   tokenString,
		Expires: time.Now().AddDate(0, 1 /* month */, 0),
	})
	ctx.RenderResponse(http.StatusOK, responses.NewSuccessResponse(nil))
}
