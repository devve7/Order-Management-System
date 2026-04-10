package order

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, customerID CustomerID, status OrderStatus) (OrderID, error)

	Get(ctx context.Context, id OrderID) (*Order, error)

	Update(ctx context.Context, order *Order) error

	GetAll(ctx context.Context) ([]*Order, error)
}
