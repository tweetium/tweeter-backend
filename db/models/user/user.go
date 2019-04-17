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

// ErrPasswordTooShort is the error returned by Create when the password is too short
var ErrPasswordTooShort = errors.New("password too short")

// ErrUserEmailAlreadyExists is the error returned by Create when the email already exists
var ErrUserEmailAlreadyExists = errors.New("user already exists for email")

// ErrUserIDNotValid is the error when trying to parse a userID string that is not a positive integer
var ErrUserIDNotValid = errors.New("user id is not a valid positive integer")

// ID is a type alias for the user's ID type
type ID uint64

// ParseID parses a string into a ID
func ParseID(idString string) (ID, error) {
	// Parse the string as base10 into a uint64
	idInt, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		return 0, ErrUserIDNotValid
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

		logrus.WithError(err).Warn("Failed to insert user into DB")
		return User{}, ErrInternalError
	}

	user, err = Find(FindByID{ID: id})
	if err != nil {
		return User{}, ErrInternalError
	}
	return user, nil
}

// Find finds and returns a user if found, otherwise err
func Find(f FindBy) (user User, err error) {
	user, err = f.Find()
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserNotFound
		}

		logrus.WithError(err).Warn("Failed to find user from DB")
		return user, ErrInternalError
	}

	return
}

// ComparePassword returns true if the password matches the stored password
func (user User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
