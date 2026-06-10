package security

import (
	"fmt"
	"strconv"
	"time"

	"Order-Management-System/internal/application/auth"
	domain_user "Order-Management-System/internal/domain/user"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret    []byte
	accessTTL time.Duration
	issuer    string
}

func NewJWTManager(secret string, accessTTL time.Duration, issuer string) (*JWTManager, error) {
	if secret == "" {
		return nil, ErrEmptyJWTSecret
	}
	if accessTTL <= 0 {
		return nil, ErrInvalidJWTAccessTTL
	}
	if issuer == "" {
		return nil, ErrEmptyJWTIssuer
	}
	return &JWTManager{
		secret:    []byte(secret),
		accessTTL: accessTTL,
		issuer:    issuer,
	}, nil
}

// type TokenManager interface {
// 	CreateAccessToken(actor Actor) (AccessToken, error)
// 	ParseAccessToken(token AccessToken) (Actor, error)
// }

func (m *JWTManager) CreateAccessToken(actor auth.Actor) (auth.AccessToken, error) {
	now := time.Now()

	claims := jwtClaims{
		Role: actor.Role().String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(int64(actor.ID()), 10),
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	accessToken, err := auth.NewAccessToken(signedToken)
	if err != nil {
		return "", fmt.Errorf("create access token: %w", err)
	}

	return accessToken, nil
}

func (m *JWTManager) ParseAccessToken(token auth.AccessToken) (auth.Actor, error) {
	claims := &jwtClaims{}

	parsedToken, err := jwt.ParseWithClaims(
		token.String(),
		claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, auth.ErrInvalidToken
			}

			return m.secret, nil
		},
		jwt.WithIssuer(m.issuer),
	)
	if err != nil {
		return auth.Actor{}, auth.ErrInvalidToken
	}

	if !parsedToken.Valid {
		return auth.Actor{}, auth.ErrInvalidToken
	}

	userIDInt, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return auth.Actor{}, auth.ErrInvalidToken
	}

	userID, err := domain_user.NewUserID(userIDInt)
	if err != nil {
		return auth.Actor{}, auth.ErrInvalidToken
	}

	role, err := domain_user.NewRole(claims.Role)
	if err != nil {
		return auth.Actor{}, auth.ErrInvalidToken
	}

	actor := auth.NewActor(userID, role)

	return actor, nil
}
