package dao

import (
	"errors"
	"time"
)

var (
	// ErrInvalidCredentials error is used if a user tries
	// to login with an incorrect email address or password.
	ErrInvalidCredentials = errors.New("dao: invalid credentials")

	// ErrDuplicateEmail error is used if a user tries to signup
	// with an email address that's already in use.
	ErrDuplicateEmail = errors.New("dao: duplicate email")
)

type User struct {
	Id       int64
	Email    string
	Password []byte
	Name     string
	Created  time.Time
	Active   bool
}
