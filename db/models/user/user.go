package user

import (
	"strconv"
	"time"

	"tweeter/db"
)

// User represents a user of the application
type User struct {
	ID       ID `db:"user_id"`
	Email    string
	Password string
	Created  time.Time
	Modified time.Time
}

// ID is a type alias for the user's ID type
type ID uint64

// ParseID parses a string into a ID
func ParseID(idString string) (ID, error) {
	idInt, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return 0, err
	}

	id := ID(idInt)
	return id, nil
}

// Create creates a user in the database and returns the model
func Create(email, password string) (user User, err error) {
	var id ID
	err = db.DB.QueryRowx("INSERT INTO users (email, password) VALUES ($1, $2) RETURNING user_id", email, password).Scan(&id)
	if err != nil {
		return
	}

	user, err = Get(id)
	return
}

// Get queries and return a user if found, otherwise err
func Get(userID ID) (user User, err error) {
	err = db.DB.QueryRowx("SELECT * FROM users WHERE user_id = $1", userID).StructScan(&user)
	return
}
