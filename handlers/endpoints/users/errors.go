package users

import (
	"fmt"
	"tweeter/db/models/user"
	"tweeter/handlers/responses"
)

// ErrEmailAlreadyExists is the responses.Error returned when email already exists
var ErrEmailAlreadyExists = func(email string) responses.Error {
	return responses.Error{
		Title:  "Email Already Exists",
		Detail: fmt.Sprintf("User already exists for %s", email),
	}
}

// ErrTooShortPassword is the responses.Error returned when password is too short
var ErrTooShortPassword = responses.Error{
	Title:  "Password Too Short",
	Detail: fmt.Sprintf("Password is too short, minimum password length: %d", user.MinimumPasswordLength),
}
