// Package product ...
package product

import (
	"Order-Management-System/internal/application/ports"
	domain_product "Order-Management-System/internal/domain/product"
	"context"
)

type UseCase struct {
	repo domain_product.Repository
}

func NewUseCase(repo domain_product.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) CreateProduct(
	ctx context.Context,
	name string,
	price int64,
	stock int64,
) (int64, error) {
	productName, err := domain_product.NewProductName(name)
	if err != nil {
		return 0, err
	}
	productPrice, err := domain_product.NewPrice(price)
	if err != nil {
		return 0, err
	}
	productStock, err := domain_product.NewStock(stock)
	if err != nil {
		return 0, err
	}
	product := domain_product.NewProduct(productName, productPrice, productStock)
	id, err := u.repo.Create(ctx, product)
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func (u *UseCase) GetProduct(ctx context.Context, productID int64) (*ProductDTO, error) {
	id, err := domain_product.NewProductID(productID)
	if err != nil {
		return nil, err
	}
	product, err := u.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &ProductDTO{
		ID:     int64(product.ID()),
		Name:   string(product.Name()),
		Price:  int64(product.Price()),
		Stock:  int64(product.Stock()),
		Active: product.Active(),
	}, nil
}

func (u *UseCase) GetProducts(ctx context.Context) ([]*ProductDTO, error) {
	products, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	productsDTO := make([]*ProductDTO, 0, len(products))
	for _, product := range products {
		productsDTO = append(productsDTO, &ProductDTO{
			ID:     int64(product.ID()),
			Name:   string(product.Name()),
			Price:  int64(product.Price()),
			Stock:  int64(product.Stock()),
			Active: product.Active(),
		})
	}
	return productsDTO, nil
}

func (u *UseCase) ChangePrice(ctx context.Context, id int64, price int64) error {
	productID, err := domain_product.NewProductID(id)
	if err != nil {
		return err
	}
	productPrice, err := domain_product.NewPrice(price)
	if err != nil {
		return err
	}
	product, err := u.repo.Get(ctx, productID)
	if err != nil {
		return err
	}
	product.ChangePrice(productPrice)
	err = u.repo.Update(ctx, product)
	return err
}

func (u *UseCase) AddStock(ctx context.Context, id int64, amount int64) error {
	productID, err := domain_product.NewProductID(id)
	if err != nil {
		return err
	}
	productStock, err := domain_product.NewStock(amount)
	if err != nil {
		return err
	}
	product, err := u.repo.Get(ctx, productID)
	if err != nil {
		return err
	}
	product.AddStock(productStock)
	err = u.repo.Update(ctx, product)
	return err
}

func (u *UseCase) RemoveStock(ctx context.Context, id int64, amount int64) error {
	productID, err := domain_product.NewProductID(id)
	if err != nil {
		return err
	}
	productStock, err := domain_product.NewStock(amount)
	if err != nil {
		return err
	}
	product, err := u.repo.Get(ctx, productID)
	if err != nil {
		return err
	}
	err = product.RemoveStock(productStock)
	if err != nil {
		return err
	}
	err = u.repo.Update(ctx, product)
	return err
}

func (u *UseCase) DeactivateProduct(ctx context.Context, id int64) error {
	productID, err := domain_product.NewProductID(id)
	if err != nil {
		return err
	}
	product, err := u.repo.Get(ctx, productID)
	if err != nil {
		return err
	}

	err = product.Deactivate()
	if err != nil {
		return err
	}

	err = u.repo.Update(ctx, product)
	return err
}

func (u *UseCase) ActivateProduct(ctx context.Context, id int64) error {
	productID, err := domain_product.NewProductID(id)
	if err != nil {
		return err
	}
	product, err := u.repo.Get(ctx, productID)
	if err != nil {
		return err
	}

	err = product.Activate()
	if err != nil {
		return err
	}

	err = u.repo.Update(ctx, product)
	return err
}

func (u *UseCase) GetSnapshot(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
	id, err := domain_product.NewProductID(productID)
	if err != nil {
		return ports.ProductSnapshot{}, err
	}
	product, err := u.repo.Get(ctx, id)
	if err != nil {
		return ports.ProductSnapshot{}, err
	}
	productSnapshot := ports.ProductSnapshot{
		ProductID: int64(product.ID()),
		Name:      string(product.Name()),
		Price:     int64(product.Price()),
	}

	return productSnapshot, nil
}

func (u *UseCase) EnsureAvailable(ctx context.Context, productID int64, quantity int64) error {
	id, err := domain_product.NewProductID(productID)
	if err != nil {
		return err
	}

	qty, err := domain_product.NewStock(quantity)
	if err != nil {
		return err
	}

	product, err := u.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	return product.EnsureAvailable(qty)
}
