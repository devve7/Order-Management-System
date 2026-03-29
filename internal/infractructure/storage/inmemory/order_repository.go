package inmemory

import (
	do "Order-Management-System/internal/domain/order"
	"sync"
)

type InMemoryOrderRepository struct {
	nextID do.OrderID
	orders map[do.OrderID]*do.Order
	mtx    sync.RWMutex
}

func NewInmemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		nextID: 1,
		orders: make(map[do.OrderID]*do.Order),
	}
}

func (r *InMemoryOrderRepository) Save(order *do.Order) error {
	if order == nil {
		return do.ErrOrderEmpty
	}
	orderID := order.ID()
	if order.ID() == 0 {
		return do.ErrInvalidID
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.orders[orderID] = order

	return nil
}

func (r *InMemoryOrderRepository) Update(order *do.Order) error {
	if order == nil {
		return do.ErrOrderEmpty
	}
	orderID := order.ID()
	if _, ok := r.orders[orderID]; !ok {
		return do.ErrOrderNotFound
	}
	if order.ID() == 0 {
		return do.ErrInvalidID
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.orders[orderID] = order

	return nil
}

func (r *InMemoryOrderRepository) Get(id do.OrderID) (*do.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return nil, do.ErrOrderNotFound
	}
	return order.Clone(), nil
}

func (r *InMemoryOrderRepository) NextID() (do.OrderID, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	result := r.nextID
	r.nextID += 1
	return result, nil
}
