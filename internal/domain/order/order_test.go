package order

import (
	"errors"
	"testing"
	"time"
)

func errCompare(t *testing.T, err error, expectedErr error) {
	t.Helper()
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected (%v), got (%v)", expectedErr, err)
	}
}

var (
	orderID    OrderID    = 10
	customerID CustomerID = 10
)

func TestAddItem(t *testing.T) {
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
				order := NewOrder(orderID, customerID, "created", time.Now())
				return order
			},
			expectedErr: nil,
		},
		{
			name: "add nil item",
			item: nil,
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
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
				order := NewOrder(orderID, customerID, "created", time.Now())
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
				order := NewOrder(orderID, customerID, "created", time.Now())
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
			errCompare(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				if len(order.items) != 1 {
					t.Errorf("item was not added")
				}
			}
		})
	}
}

func TestRemoveItem(t *testing.T) {
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
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.AddItem(item)
				return order
			},
			expectedErr: nil,
		},
		{
			name: "delete from empty",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				return order
			},
			expectedErr: ErrOrderItemNotFound,
		},
		{
			name: "delete unknown",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
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
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.status = StatusCancelled
				return order
			},
			expectedErr: ErrOrderCancelled,
		},
		{
			name: "delete from shipped",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
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
			errCompare(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				if len(order.items) != 0 {
					t.Errorf("item was not removed")
				}
			}
		})
	}
}

func TestPay(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *Order
		expectedErr error
	}{
		{
			name: "ok",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				item := &OrderItem{
					id:    2,
					name:  "item",
					price: 100,
				}
				order.AddItem(item)
				return order
			},
			expectedErr: nil,
		},
		{
			name: "pay empty order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				return order
			},
			expectedErr: ErrOrderEmpty,
		},
		{
			name: "pay paid order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				item := &OrderItem{
					id:    2,
					name:  "item",
					price: 100,
				}
				order.AddItem(item)
				order.status = StatusPaid
				return order
			},
			expectedErr: ErrCannotPay,
		},
		{
			name: "pay shipped order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				item := &OrderItem{
					id:    2,
					name:  "item",
					price: 100,
				}
				order.AddItem(item)
				order.status = StatusShipped
				return order
			},
			expectedErr: ErrCannotPay,
		},
		{
			name: "pay cancelled order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				item := &OrderItem{
					id:    2,
					name:  "item",
					price: 100,
				}
				order.AddItem(item)
				order.status = StatusCancelled
				return order
			},
			expectedErr: ErrCannotPay,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setup()
			err := order.Pay()
			errCompare(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				if order.status != StatusPaid {
					t.Errorf("item was not paid")
				}
			}
		})
	}
}

func TestShip(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *Order
		expectedErr error
	}{
		{
			name: "ok",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.status = StatusPaid
				return order
			},
			expectedErr: nil,
		},
		{
			name: "ship shipped order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.status = StatusShipped
				return order
			},
			expectedErr: ErrCannotShip,
		},
		{
			name: "ship cancelled order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.status = StatusCancelled
				return order
			},
			expectedErr: ErrCannotShip,
		},
		{
			name: "ship created order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.status = StatusCreated
				return order
			},
			expectedErr: ErrCannotShip,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setup()
			err := order.Ship()
			errCompare(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				if order.status != StatusShipped {
					t.Errorf("item was not shipped")
				}
			}
		})
	}
}

func TestCancel(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *Order
		expectedErr error
	}{
		{
			name: "ok",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.status = StatusCreated
				return order
			},
			expectedErr: nil,
		},
		{
			name: "cancel shipped order",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				order.status = StatusShipped
				return order
			},
			expectedErr: ErrCannotCancel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setup()
			err := order.Cancel()
			errCompare(t, err, tt.expectedErr)
			if tt.expectedErr == nil {
				if order.status != StatusCancelled {
					t.Errorf("item was not cancelled")
				}
			}
		})
	}
}

func TestTotal(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Order
		expected Price
	}{
		{
			name: "ok",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				item := &OrderItem{
					id:    2,
					name:  "item",
					price: 100,
				}
				order.AddItem(item)
				return order
			},
			expected: 100,
		},
		{
			name: "total of many",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				item := &OrderItem{
					id:    1,
					name:  "item",
					price: 100,
				}
				item2 := &OrderItem{
					id:    2,
					name:  "item",
					price: 100,
				}
				order.AddItem(item)
				order.AddItem(item2)
				return order
			},
			expected: 200,
		},
		{
			name: "ok",
			setup: func() *Order {
				order := NewOrder(orderID, customerID, "created", time.Now())
				return order
			},
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := tt.setup()
			price := order.Total()
			if price != tt.expected {
				t.Errorf("expected (%v), got (%v)", tt.expected, price)
			}
		})
	}
}
