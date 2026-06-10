// Package security ...
package security

import (
	domain_user "Order-Management-System/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

type BcryptPasswordHasher struct{}

func (h *BcryptPasswordHasher) Hash(password string) (domain_user.PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	passwordHash, err := domain_user.NewPasswordHash(string(hash))
	if err != nil {
		return "", err
	}
	return passwordHash, nil
}

func (h *BcryptPasswordHasher) Compare(password string, hash domain_user.PasswordHash) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
}
