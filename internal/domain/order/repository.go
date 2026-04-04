package order

import "context"

type Repository interface {
	Save(ctx context.Context, order *Order) error

	Get(ctx context.Context, id OrderID) (*Order, error)

	Update(ctx context.Context, order *Order) error

	NextID(ctx context.Context) (OrderID, error)

	GetAll(ctx context.Context) ([]*Order, error)
}
