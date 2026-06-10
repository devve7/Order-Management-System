// Package user ...
package user

import (
	"net/mail"
	"strings"
)

type UserID int64

func NewUserID(id int64) (UserID, error) {
	if id <= 0 {
		return 0, ErrInvalidUserID
	}
	return UserID(id), nil
}

type Email string

func NewEmail(input string) (Email, error) {
	if input == "" {
		return "", ErrInvalidEmail
	}
	if len(input) > 255 {
		return "", ErrInvalidEmail
	}
	toLower := strings.ToLower(input)
	email, err := mail.ParseAddress(toLower)
	if err != nil {
		return "", ErrInvalidEmail
	}
	return Email(email.Address), nil
}

type PasswordHash string

func NewPasswordHash(hash string) (PasswordHash, error) {
	if hash == "" {
		return "", ErrEmptyPasswordHash
	}

	return PasswordHash(hash), nil
}
