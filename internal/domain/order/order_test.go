package order

import (
	"errors"
	"testing"
)

func TestAddItem(t *testing.T) {
	orderID, _ := NewOrderID("orderID")
	CustomerID, _ := NewCustomerID("customerID")
	tests := []struct {
		name        string
		item        *OrderItem
		setup       func() *Order
		expectedErr error
	}{
		{
			name: "ok",
			item: &OrderItem{
				id:    1,
				name:  "hello",
				price: 100,
			},
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				return order
			},
			expectedErr: nil,
		},
		{
			name: "add nil item",
			item: nil,
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				return order
			},
			expectedErr: ErrOrderItemNotFound,
		},
		{
			name: "add for cancelled",
			item: &OrderItem{
				id:    1,
				name:  "hello",
				price: 100,
			},
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				order.status = StatusCancelled
				return order
			},
			expectedErr: ErrOrderCancelled,
		},
		{
			name: "add for shipped",
			item: &OrderItem{
				id:    1,
				name:  "hello",
				price: 100,
			},
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				order.status = StatusShipped
				return order
			},
			expectedErr: ErrOrderShipped,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setup()
			err := order.AddItem(tt.item)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
		})
	}
}

func TestRemoveItem(t *testing.T) {
	orderID, _ := NewOrderID("orderID")
	CustomerID, _ := NewCustomerID("customerID")
	item := &OrderItem{
		id:    1,
		name:  "hello",
		price: 100,
	}
	itemID := item.ID()
	tests := []struct {
		name        string
		setup       func() *Order
		expectedErr error
	}{
		{
			name: "delete first time",
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				order.AddItem(item)
				return order
			},
			expectedErr: nil,
		},
		{
			name: "delete from empty",
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				return order
			},
			expectedErr: ErrOrderItemNotFound,
		},
		{
			name: "delete unknown",
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				otherItem := &OrderItem{
					id:    2,
					name:  "other",
					price: 100,
				}
				order.AddItem(otherItem)
				return order
			},
			expectedErr: ErrOrderItemNotFound,
		},
		{
			name: "delete from cancelled",
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				order.status = StatusCancelled
				return order
			},
			expectedErr: ErrOrderCancelled,
		},
		{
			name: "delete from shipped",
			setup: func() *Order {
				order := NewOrder(orderID, CustomerID)
				order.status = StatusShipped
				return order
			},
			expectedErr: ErrOrderShipped,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setup()
			err := order.RemoveItem(itemID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
		})
	}
}
