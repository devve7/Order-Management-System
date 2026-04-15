// Package ports
package ports

import "context"

type ProductSnapshot struct {
	ProductID int64
	Name      string
	Price     int64
}

type ProductForOrder interface {
	GetSnapshot(ctx context.Context, productID int64) (ProductSnapshot, error)
	EnsureAvailable(ctx context.Context, productID int64, quantity int64) error
}
