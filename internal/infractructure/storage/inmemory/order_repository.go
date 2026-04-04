// Package inmemory ...
package inmemory

import (
	domain_order "Order-Management-System/internal/domain/order"
	"context"
	"sync"
)

type InMemoryOrderRepository struct {
	nextID domain_order.OrderID
	orders map[domain_order.OrderID]*domain_order.Order
	mtx    sync.RWMutex
}

func NewInmemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		nextID: 1,
		orders: make(map[domain_order.OrderID]*domain_order.Order),
	}
}

func (r *InMemoryOrderRepository) Save(ctx context.Context, order *domain_order.Order) error {
	if order == nil {
		return domain_order.ErrOrderEmpty
	}
	orderID := order.ID()
	if order.ID() == 0 {
		return domain_order.ErrInvalidID
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.orders[orderID] = order

	return nil
}

func (r *InMemoryOrderRepository) Update(ctx context.Context, order *domain_order.Order) error {
	if order == nil {
		return domain_order.ErrOrderEmpty
	}
	orderID := order.ID()
	if _, ok := r.orders[orderID]; !ok {
		return domain_order.ErrOrderNotFound
	}
	if order.ID() == 0 {
		return domain_order.ErrInvalidID
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.orders[orderID] = order

	return nil
}

func (r *InMemoryOrderRepository) NextID(ctx context.Context) (domain_order.OrderID, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	result := r.nextID
	r.nextID += 1
	return result, nil
}

func (r *InMemoryOrderRepository) Get(ctx context.Context, id domain_order.OrderID) (*domain_order.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return nil, domain_order.ErrOrderNotFound
	}
	return order.Clone(), nil
}

func (r *InMemoryOrderRepository) GetAll(ctx context.Context) ([]*domain_order.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	result := make([]*domain_order.Order, 0, len(r.orders))
	for _, v := range r.orders {
		result = append(result, v.Clone())
	}
	return result, nil
}
