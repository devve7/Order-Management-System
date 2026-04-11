// Package inmemory
package inmemory

import (
	"context"
	"sync"

	app_product "Order-Management-System/internal/application/product"
	domain_order "Order-Management-System/internal/domain/order"
)

type ProductRecord struct {
	ID        domain_order.ProductID
	Name      domain_order.ProductName
	Price     domain_order.Price
	Available domain_order.Quantity
}

type InMemoryService struct {
	mtx      sync.RWMutex
	products map[domain_order.ProductID]ProductRecord
}

func NewInMemoryService() *InMemoryService {
	return &InMemoryService{
		products: make(map[domain_order.ProductID]ProductRecord),
	}
}

func (s *InMemoryService) AddProduct(
	id domain_order.ProductID,
	name domain_order.ProductName,
	price domain_order.Price,
	available domain_order.Quantity,
) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.products[id] = ProductRecord{
		ID:        id,
		Name:      name,
		Price:     price,
		Available: available,
	}
}

func (s *InMemoryService) GetSnapshot(
	ctx context.Context,
	productID domain_order.ProductID,
) (app_product.Snapshot, error) {

	s.mtx.RLock()
	defer s.mtx.RUnlock()

	product, ok := s.products[productID]
	if !ok {
		return app_product.Snapshot{}, app_product.ErrProductNotFound
	}

	return app_product.Snapshot{
		ProductID: product.ID,
		Name:      product.Name,
		Price:     product.Price,
	}, nil
}

func (s *InMemoryService) EnsureAvailable(
	ctx context.Context,
	productID domain_order.ProductID,
	quantity domain_order.Quantity,
) error {

	if quantity <= 0 {
		return domain_order.ErrInvalidQuantity
	}

	s.mtx.RLock()
	defer s.mtx.RUnlock()

	product, ok := s.products[productID]
	if !ok {
		return app_product.ErrProductNotFound
	}

	if product.Available < quantity {
		return app_product.ErrInsufficientStock
	}

	return nil
}
