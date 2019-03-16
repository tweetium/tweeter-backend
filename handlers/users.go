package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"tweeter/db/models/user"
	"tweeter/handlers/responses"
)

// UsersHandler is the /api/v1/users/ endpoint
func UsersHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPut {
		handleUserCreate(w, req)
	} else {
		renderErrors(w, http.StatusBadRequest, responses.Error{
			Title: "Unsupported Method", Detail: fmt.Sprintf("Method %s is not supported (yet)", req.Method),
		})
	}
}

func parseID(w http.ResponseWriter, req *http.Request) (ID user.ID, ok bool) {
	idString := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/api/v1/users/"), "/")
	id, err := user.ParseID(idString)
	if err != nil {
		renderErrors(w, http.StatusBadRequest, responses.Error{
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
		log.Printf("Failed to read request body, err: %s", err)
		renderErrors(w, http.StatusBadRequest, responses.Error{
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
		renderErrors(w, http.StatusBadRequest, responses.Error{
			Title: "Invalid Body", Detail: fmt.Sprintf("Failed to parse request body as json, err: %s", err),
		})
		return
	}

	// TODO: classify this into internal error (couldn't connect / exec on db)
	// or other errors like invalid user (email already exists, etc)
	var newUser user.User
	newUser, err = user.Create(createReq.Email, createReq.Password)
	if err != nil {
		renderErrors(w, http.StatusInternalServerError, responses.Error{
			Title: "User Creation Error", Detail: fmt.Sprintf("Failed to create user, err: %s", err),
		})
		return
	}

	render(w, http.StatusOK, responses.NewSuccessResponse(struct {
		ID user.ID
	}{
		newUser.ID,
	}))
}
