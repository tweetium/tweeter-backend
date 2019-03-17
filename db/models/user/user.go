package user

import (
	"database/sql"
	"errors"
	"strconv"
	"time"

	"tweeter/db"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user of the application
type User struct {
	ID       ID `db:"user_id"`
	Email    string
	Password string
	Created  time.Time
	Modified time.Time
}

// MinimumPasswordLength is the minimum length of a password
var MinimumPasswordLength = 6

// ErrUserNotFound is the error returned when the given user does not exist
var ErrUserNotFound = errors.New("could not find user from information given")

// ErrInternalError is the error returned when normal flow failed (db error, timeout, etc)
var ErrInternalError = errors.New("internal error, retry later")

// ErrMismatchedPassword is the error returned by Validate when password do not match stored password
var ErrMismatchedPassword = errors.New("invalid password given")

// ErrPasswordTooShort is the error returned by Create when the password is too short
var ErrPasswordTooShort = errors.New("password too short")

// ErrUserEmailAlreadyExists is the error returned by Create when the email already exists
var ErrUserEmailAlreadyExists = errors.New("user already exists for email")

// ErrUserIDNotInterger is the error when trying to parse a userID string that is not an integer
var ErrUserIDNotInterger = errors.New("user id is not an integer")

// ID is a type alias for the user's ID type
type ID uint64

// ParseID parses a string into a ID
func ParseID(idString string) (ID, error) {
	// Parse the string as base10 into a int64
	idInt, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return 0, ErrUserIDNotInterger
	}

	id := ID(idInt)
	return id, nil
}

// Create creates a user in the database and returns the model
func Create(email, password string) (user User, err error) {
	if len(password) < MinimumPasswordLength {
		return User{}, ErrPasswordTooShort
	}

	var hashedPasswordBytes []byte
	passwordBytes := []byte(password)
	hashedPasswordBytes, err = bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return User{}, ErrInternalError
	}

	var id ID
	err = db.DB.QueryRowx(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING user_id",
		email, string(hashedPasswordBytes),
	).Scan(&id)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return User{}, ErrUserEmailAlreadyExists
		}

		logrus.WithFields(logrus.Fields{"err": err}).Warn("Failed to insert user into DB")
		return User{}, ErrInternalError
	}

	user, err = Get(id)
	return
}

// Get queries and return a user if found, otherwise err
func Get(userID ID) (user User, err error) {
	err = db.DB.QueryRowx("SELECT * FROM users WHERE user_id = $1", userID).StructScan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserNotFound
		}

		return user, ErrInternalError
	}

	return
}

// Validate returns nil if email / password match, otherwise err
func Validate(email, password string) (err error) {
	var hashedPassword string
	err = db.DB.QueryRowx("SELECT password FROM users WHERE email = $1", email).Scan(&hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}

		return ErrInternalError
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return ErrMismatchedPassword
	}

	return nil
}
