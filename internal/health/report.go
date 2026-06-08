// Package health ...
package health

import "context"

type HealthChecker interface {
	Check(ctx context.Context) HealthReport
}

type HealthReport struct {
	Status string                 `json:"status"`
	Checks map[string]CheckStatus `json:"checks"`
}

type CheckStatus struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}
