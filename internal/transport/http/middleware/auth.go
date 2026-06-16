package middleware

import (
	"net/http"
	"strings"

	"Order-Management-System/internal/application/auth"
)

type AuthMiddleware struct {
	tokenManager auth.TokenManager
}

func NewAuthMiddleware(tokenManager auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: tokenManager,
	}
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(authorization, bearerPrefix) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tokenValue := strings.TrimSpace(strings.TrimPrefix(authorization, bearerPrefix))
		if tokenValue == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken, err := auth.NewAccessToken(tokenValue)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		actor, err := m.tokenManager.ParseAccessToken(accessToken)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := WithActor(r.Context(), actor)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
