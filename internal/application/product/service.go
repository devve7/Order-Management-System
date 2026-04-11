// Package product ...
package product

import (
	domain_order "Order-Management-System/internal/domain/order"
	"context"
)

type Snapshot struct {
	ProductID domain_order.ProductID
	Name      domain_order.ProductName
	Price     domain_order.Price
}

type ProductService interface {
	GetSnapshot(ctx context.Context, productID domain_order.ProductID) (Snapshot, error)
	EnsureAvailable(ctx context.Context, productID domain_order.ProductID, quantity domain_order.Quantity) error
}
