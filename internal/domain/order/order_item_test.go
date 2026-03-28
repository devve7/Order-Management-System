package order

import (
	"errors"
	"testing"
)

type TestCatalog struct{}

func (c *TestCatalog) GetProduct(name string) (Product, error) {
	if name == "test" {
		return Product{
			ID:    1,
			Name:  "test",
			Price: 100,
		}, nil
	} else {
		return Product{}, ErrOrderItemNotFound
	}
}

func TestOrderItemFactory(t *testing.T) {
	var catalog TestCatalog
	factory := NewOrderItemFactory(&catalog)
	tests := []struct {
		name          string
		orderItemName string
		expectedErr   error
	}{
		{
			name:          "ok",
			orderItemName: "test",
			expectedErr:   nil,
		},
		{
			name:          "not found",
			orderItemName: "unknonw",
			expectedErr:   ErrOrderItemNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := factory.New(tt.orderItemName)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
		})
	}
}
