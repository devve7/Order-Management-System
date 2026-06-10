package user

import "errors"

var (
	ErrInvalidUserID     = errors.New("invalid user id")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrInvalidRole       = errors.New("invalid role")
	ErrEmptyPasswordHash = errors.New("empty password hash")
	ErrUserNotFound      = errors.New("user not found")
)
