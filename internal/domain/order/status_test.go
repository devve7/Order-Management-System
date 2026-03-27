package order

import (
	"errors"
	"testing"
)

func TestStatus(t *testing.T) {
	tests := []struct {
		name   string
		status string
		err    error
	}{
		{
			name:   "Created",
			status: "created",
			err:    nil,
		},
		{
			name:   "Paid",
			status: "paid",
			err:    nil,
		},
		{
			name:   "Shipped",
			status: "shipped",
			err:    nil,
		},
		{
			name:   "Cancelled",
			status: "cancelled",
			err:    nil,
		},
		{
			name:   "Unknown",
			status: "unknown",
			err:    ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := NewOrderStatus(tt.status)
			if !errors.Is(err, tt.err) {
				t.Errorf("Incorrect Status, expected (%v, %v), got (%v, %v)", tt.err, status, err, status)
			}
		})
	}
}
