package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"tweeter/db/models/user"
	"tweeter/handlers/endpoints"
	"tweeter/handlers/render"
	"tweeter/handlers/responses"
)

// Endpoint is the /api/v1/users/ endpoint that handles user CRUD apis
var Endpoint = endpoints.Endpoint{
	URL:     "/api/v1/users",
	Handler: handler,
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		handleUserCreate(w, req)
	} else {
		render.ErrorResponse(w, http.StatusBadRequest, responses.Error{
			Title: "Unsupported Method", Detail: fmt.Sprintf("Method %s is not supported (yet)", req.Method),
		})
	}
}

// parseID is unused, but will be used when user GET is added
func parseID(w http.ResponseWriter, req *http.Request) (ID user.ID, ok bool) { //nolint
	idString := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/api/v1/users/"), "/")
	id, err := user.ParseID(idString)
	if err != nil {
		render.ErrorResponse(w, http.StatusBadRequest, responses.Error{
			Title: "Invalid User ID", Detail: fmt.Sprintf("Failed to parse ID from %s, err: %s", idString, err),
		})
		return 0, false
	}

	return id, true
}

func handleUserCreate(w http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		// This is unexpected (but possible), so let's log this internally here
		logrus.WithFields(logrus.Fields{"err": err}).Warn("Failed to read request body")
		render.ErrorResponse(w, http.StatusBadRequest, responses.Error{
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
		render.ErrorResponse(w, http.StatusBadRequest, responses.Error{
			Title: "Invalid Body", Detail: fmt.Sprintf("Failed to parse request body as json, err: %s", err),
		})
		return
	}

	var newUser user.User
	newUser, err = user.Create(createReq.Email, createReq.Password)
	if err != nil {
		switch err {
		case user.ErrInternalError:
			render.ErrorResponse(w, http.StatusInternalServerError, responses.ErrInternalError)
		case user.ErrPasswordTooShort:
			render.ErrorResponse(w, http.StatusBadRequest, ErrPasswordTooShort)
		case user.ErrUserEmailAlreadyExists:
			render.ErrorResponse(w, http.StatusBadRequest, ErrEmailAlreadyExists(createReq.Email))
		default:
			// Logged as error because this indicates a programmer error, should fix the code if this happens
			logrus.WithFields(logrus.Fields{"err": err}).Error("Uncaught error for user.Create")
			render.ErrorResponse(w, http.StatusInternalServerError, responses.ErrInternalError)
		}
		return
	}

	render.Response(w, http.StatusOK, responses.NewSuccessResponse(struct {
		ID user.ID
	}{
		newUser.ID,
	}))
}
