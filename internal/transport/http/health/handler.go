// Package health ...
package health

import (
	health "Order-Management-System/internal/health"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

type HealthHandler struct {
	checker health.HealthChecker
	logger  *logrus.Logger
}

func NewHealthHandler(checker health.HealthChecker) *HealthHandler {
	return &HealthHandler{
		checker: checker,
	}
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	report := h.checker.Check(ctx)

	var statusCode int
	switch report.Status {
	case "ok", "degraded":
		statusCode = http.StatusOK
	case "unavailable":
		statusCode = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(report); err != nil {
		h.logger.WithError(err).Error("failed to write health check response")
	}
}
