package order

import (
	"errors"
	"testing"
)

func TestOrderID(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		err  error
	}{
		{
			name: "ok",
			id:   10,
			err:  nil,
		},
		{
			name: "negative",
			id:   -10,
			err:  ErrInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orderID, err := NewOrderID(tt.id)
			if !errors.Is(err, tt.err) {
				t.Errorf("expected (%v), got (%v)", tt.err, err)
			}
			if tt.err == nil {
				if orderID != OrderID(tt.id) {
					t.Errorf("expected (%v), got (%v)", tt.id, orderID)
				}
			}
		})
	}
}

func TestCustomerID(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		err  error
	}{
		{
			name: "ok",
			id:   10,
			err:  nil,
		},
		{
			name: "negative",
			id:   -10,
			err:  ErrInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerID, err := NewCustomerID(tt.id)
			if !errors.Is(err, tt.err) {
				t.Errorf("expected (%v), got (%v)", tt.err, err)
			}
			if tt.err == nil {
				if customerID != CustomerID(tt.id) {
					t.Errorf("expected (%v), got (%v)", tt.id, customerID)
				}
			}
		})
	}
}

func TestItemID(t *testing.T) {
	tests := []struct {
		name string
		id   int64
		err  error
	}{
		{
			name: "ok",
			id:   10,
			err:  nil,
		},
		{
			name: "negative",
			id:   -10,
			err:  ErrInvalidID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			itemID, err := NewItemID(tt.id)
			if !errors.Is(err, tt.err) {
				t.Errorf("expected (%v), got (%v)", tt.err, err)
			}
			if tt.err == nil {
				if itemID != ItemID(tt.id) {
					t.Errorf("expected (%v), got (%v)", tt.id, itemID)
				}
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
