package security

import (
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}
