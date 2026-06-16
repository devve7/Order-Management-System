// Package auth ...
package auth

import (
	"Order-Management-System/internal/application/auth"
	"Order-Management-System/internal/domain/user"
	"Order-Management-System/internal/transport/http/middleware"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	authService *auth.AuthService
	logger      *logrus.Logger
}

func NewAuthHandler(authService *auth.AuthService, logger *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func mapError(err error) int {
	switch {
	case errors.Is(err, user.ErrInvalidEmail),
		errors.Is(err, user.ErrInvalidUserID):
		return http.StatusBadRequest

	case errors.Is(err, auth.ErrInvalidCredentials),
		errors.Is(err, auth.ErrInvalidToken):
		return http.StatusUnauthorized

	case errors.Is(err, user.ErrUserAlreadyExists):
		return http.StatusConflict

	case errors.Is(err, user.ErrUserNotFound):
		return http.StatusNotFound

	default:
		return http.StatusInternalServerError
	}
}

func (h *AuthHandler) writeError(w http.ResponseWriter, r *http.Request, err error, status int) {
	entry := h.logger.WithFields(logrus.Fields{
		"method":      r.Method,
		"uri":         r.RequestURI,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"status":      status,
	}).WithError(err)

	if status >= 500 {
		entry.Error("request failed")
	} else {
		entry.Warn("request failed")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: err.Error(),
	}

	if status == 500 {
		resp.Error = "internal server error"
	}

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		h.logger.WithError(err).Error("failed to write error response")
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		h.writeError(w, r, errors.New("email is required"), http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		h.writeError(w, r, errors.New("password is required"), http.StatusBadRequest)
		return
	}

	out, err := h.authService.Register(r.Context(), auth.RegisterInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	resp := RegisterResponse{
		ID: int64(out.ID),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to write register auth response")
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	if req.Email == "" {
		h.writeError(w, r, errors.New("email is required"), http.StatusBadRequest)
		return
	}
	if req.Password == "" {
		h.writeError(w, r, errors.New("password is required"), http.StatusBadRequest)
		return
	}
	out, err := h.authService.Login(r.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	resp := LoginResponse{
		AccessToken: out.AccessToken.String(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to write login auth response")
	}
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	actor, ok := middleware.ActorFromContext(r.Context())
	if !ok {
		h.writeError(w, r, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	resp := MeResponse{
		ID:   int64(actor.ID()),
		Role: actor.Role().String(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to write me response")
	}
}
