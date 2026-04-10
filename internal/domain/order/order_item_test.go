package order

import (
	"errors"
	"testing"
)

type TestCatalog struct{}

func (c *TestCatalog) GetProduct(productID ProductID) (Product, error) {
	if productID == 1 {
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
		name               string
		orderItemProductID int
		expectedErr        error
	}{
		{
			name:               "ok",
			orderItemProductID: 1,
			expectedErr:        nil,
		},
		{
			name:               "not found",
			orderItemProductID: 2,
			expectedErr:        ErrOrderItemNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := factory.New(ProductID(tt.orderItemProductID), 1)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
		})
	}
}
