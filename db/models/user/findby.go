package user

import "tweeter/db"

// FindBy is any type that can find a user
type FindBy interface {
	Find() (User, error)
}

// FindByID is used to find a user by ID
type FindByID struct {
	ID ID
}

// Find finds a user by the id provided
func (f FindByID) Find() (user User, err error) {
	err = db.DB.QueryRowx("SELECT * FROM users WHERE user_id = $1", f.ID).StructScan(&user)
	return
}

// FindByEmail is used to find a user by Email
type FindByEmail struct {
	Email string
}

// Find finds a user by the email provided
func (f FindByEmail) Find() (user User, err error) {
	err = db.DB.QueryRowx("SELECT * FROM users WHERE email = $1", f.Email).StructScan(&user)
	return
}
