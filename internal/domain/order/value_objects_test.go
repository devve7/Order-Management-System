package order

import (
	"errors"
	"testing"
)

func TestOrderID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		err  error
	}{
		{
			name: "ok",
			id:   "1234",
			err:  nil,
		},
		{
			name: "empty id",
			id:   "",
			err:  ErrEmptyID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID, err := NewOrderID(tt.id)
			if !errors.Is(err, tt.err) || tt.id != string(orderID) {
				t.Errorf("Incorrect OrderID, expected (%v, %v), got (%v, %v)", tt.err, tt.id, err, string(orderID))
			}
		})
	}
}

func TestCustomerID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		err  error
	}{
		{
			name: "ok",
			id:   "1234",
			err:  nil,
		},
		{
			name: "empty id",
			id:   "",
			err:  ErrEmptyID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerID, err := NewCustomerID(tt.id)
			if !errors.Is(err, tt.err) || tt.id != string(customerID) {
				t.Errorf("Incorrect CustomerID, expected (%v, %v), got (%v, %v)", tt.err, tt.id, err, string(customerID))
			}
		})
	}
}

func TestPrice(t *testing.T) {
	tests := []struct {
		name           string
		amount         float64
		err            error
		expectedAmount float64
	}{
		{
			name:           "ok",
			amount:         1234,
			err:            nil,
			expectedAmount: 1234,
		},
		{
			name:           "negative amount",
			amount:         -39,
			err:            ErrInvalidPrice,
			expectedAmount: 0,
		},
		{
			name:           "zero amount",
			amount:         0,
			err:            ErrInvalidPrice,
			expectedAmount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price, err := NewPrice(tt.amount)
			if !errors.Is(err, tt.err) || float64(price) != tt.expectedAmount {
				t.Errorf("Incorrect Price, expected (%v, %v), got (%v, %v)", tt.err, tt.amount, err, tt.expectedAmount)
			}
		})
	}
}
