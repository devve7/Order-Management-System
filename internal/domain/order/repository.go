package order

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, customerID CustomerID) (OrderID, error)

	Get(ctx context.Context, id OrderID) (*Order, error)

	Update(ctx context.Context, order *Order) error

	List(ctx context.Context, params OrderListParams) ([]*Order, error)
}
