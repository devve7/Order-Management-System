package order

import "testing"

func TestOrderStatusIsValid(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{"created", StatusCreated, true},
		{"paid", StatusPaid, true},
		{"shipped", StatusShipped, true},
		{"cancelled", StatusCancelled, true},
		{"invalid", OrderStatus("unknown"), false},
		{"empty", OrderStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOrderStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    OrderStatus
		wantErr error
	}{
		{"created", "created", StatusCreated, nil},
		{"paid", "paid", StatusPaid, nil},
		{"shipped", "shipped", StatusShipped, nil},
		{"cancelled", "cancelled", StatusCancelled, nil},
		{"invalid", "other", "", ErrInvalidStatus},
		{"empty", "", "", ErrInvalidStatus},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOrderStatus(tt.input)
			if err != tt.wantErr {
				t.Fatalf("NewOrderStatus(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("NewOrderStatus(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
