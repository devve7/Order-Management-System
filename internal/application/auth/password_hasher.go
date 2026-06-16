// Package auth ...
package auth

import domain_user "Order-Management-System/internal/domain/user"

type PasswordHasher interface {
	Hash(password string) (domain_user.PasswordHash, error)
	Compare(password string, hash domain_user.PasswordHash) error
}
