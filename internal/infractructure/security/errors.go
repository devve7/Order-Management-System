package security

import "errors"

var (
	ErrEmptyJWTSecret      = errors.New("empty jwt secret")
	ErrInvalidJWTAccessTTL = errors.New("invalid jwt access ttl")
	ErrEmptyJWTIssuer      = errors.New("empty jwt issuer")
)
