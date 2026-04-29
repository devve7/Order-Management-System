// Package product ...
package product

import (
	"Order-Management-System/internal/application/ports"
	domain_product "Order-Management-System/internal/domain/product"
	"context"
	"encoding/json"
	"fmt"
)

type UseCase struct {
	repo  domain_product.Repository
	cache ports.Cache
}

func NewUseCase(repo domain_product.Repository, cache ports.Cache) *UseCase {
	return &UseCase{
		repo:  repo,
		cache: cache,
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

func (u *UseCase) GetProduct(ctx context.Context, productID int64) (ProductDTO, error) {
	id, err := domain_product.NewProductID(productID)
	if err != nil {
		return ProductDTO{}, err
	}

	cacheKey := fmt.Sprintf("product:%d", id)
	cached, err := u.cache.Get(ctx, cacheKey)
	if err == nil {
		var dto ProductDTO
		json.Unmarshal([]byte(cached), &dto)
		return dto, nil
	}

	product, err := u.repo.Get(ctx, id)
	if err != nil {
		return ProductDTO{}, err
	}

	dto := ProductDTO{
		ID:     int64(product.ID()),
		Name:   string(product.Name()),
		Price:  int64(product.Price()),
		Stock:  int64(product.Stock()),
		Active: product.Active(),
	}
	bytes, _ := json.Marshal(dto)
	u.cache.Set(ctx, cacheKey, string(bytes), 300)

	return dto, nil
}

func (u *UseCase) GetProducts(ctx context.Context) ([]ProductDTO, error) {
	cacheKey := "products:all"
	cached, err := u.cache.Get(ctx, cacheKey)
	if err == nil {
		var dtos []ProductDTO
		if err := json.Unmarshal([]byte(cached), &dtos); err == nil {
			return dtos, nil
		}
	}

	products, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]ProductDTO, 0, len(products))
	for _, p := range products {
		dtos = append(dtos, ProductDTO{
			ID:     int64(p.ID()),
			Name:   string(p.Name()),
			Price:  int64(p.Price()),
			Stock:  int64(p.Stock()),
			Active: p.Active(),
		})
	}

	if bytes, err := json.Marshal(dtos); err == nil {
		_ = u.cache.Set(ctx, cacheKey, string(bytes), 60)
	}

	return dtos, nil
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
	if err != nil {
		return err
	}

	_ = u.cache.Delete(ctx, fmt.Sprintf("product:%d", productID))
	_ = u.cache.Delete(ctx, "products:all")
	return nil
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
	if err != nil {
		return err
	}

	_ = u.cache.Delete(ctx, fmt.Sprintf("product:%d", productID))
	_ = u.cache.Delete(ctx, "products:all")
	return nil
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
	if err != nil {
		return err
	}

	_ = u.cache.Delete(ctx, fmt.Sprintf("product:%d", productID))
	_ = u.cache.Delete(ctx, "products:all")
	return nil
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
	if err != nil {
		return err
	}

	_ = u.cache.Delete(ctx, fmt.Sprintf("product:%d", productID))
	_ = u.cache.Delete(ctx, "products:all")
	return nil
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
	if err != nil {
		return err
	}

	_ = u.cache.Delete(ctx, fmt.Sprintf("product:%d", productID))
	_ = u.cache.Delete(ctx, "products:all")
	return nil
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
