// Package order ...
package order

import "time"

type Order struct {
	id         OrderID
	customerID CustomerID
	items      []OrderItem
	status     OrderStatus
	createdAt  time.Time
}

func NewOrder(id OrderID, customerID CustomerID) *Order
