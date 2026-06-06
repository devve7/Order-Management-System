package product

import "context"

type Repository interface {
	Create(ctx context.Context, product *Product) (ProductID, error)
	Get(ctx context.Context, id ProductID) (*Product, error)
	Update(ctx context.Context, product *Product) error
	List(ctx context.Context, params ProductListParams) ([]*Product, error)
}
