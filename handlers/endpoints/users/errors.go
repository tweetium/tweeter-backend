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

// ErrPasswordTooShort is the responses.Error returned when password is too short
var ErrPasswordTooShort = responses.Error{
	Title:  "Password Too Short",
	Detail: fmt.Sprintf("Password is too short, minimum password length: %d", user.MinimumPasswordLength),
}

// ErrInvalidCredentials is a deliberately generic responses.Error returned when
// anything about the credentials is incorrect. This should not give any extra
// information.
var ErrInvalidCredentials = responses.Error{
	Title:  "Invalid Credentials",
	Detail: "The credentials provided were not correct",
}

// ErrInvalidBody is the responses.Error returned when request body could not be parsed
var ErrInvalidBody = responses.Error{
	Title: "Invalid Body", Detail: "Failed to parse request body as json",
}
