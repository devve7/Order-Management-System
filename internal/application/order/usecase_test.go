package order

import (
	do "Order-Management-System/internal/domain/order"
	"errors"
	"sync"
	"testing"
)

func TestAddItem(t *testing.T) {
	customerID, _ := do.NewCustomerID(1)

	tests := []struct {
		name        string
		setup       func() *UseCase
		expectedErr error
	}{
		{
			name: "ok",
			setup: func() *UseCase {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					do.Product{
						Name: "iphone",
					},
				)
				factory := do.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				return usecase
			},
			expectedErr: nil,
		},
		{
			name: "product not found",
			setup: func() *UseCase {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					do.Product{
						Name: "iphone X",
					},
				)
				factory := do.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				return usecase
			},
			expectedErr: do.ErrProductNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup()
			orderID, _ := usecase.CreateOrder(customerID)
			itemID, err := usecase.AddItem(orderID, "iphone")
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
			if err == nil {
				orderDTO, _ := usecase.GetOrder(orderID)
				item := orderDTO.Items[0]
				if itemID != do.ItemID(item.ID) {
					t.Errorf("expected (%v), got (%v)", itemID, item.ID)
				}
			}
		})
	}
}

func TestRemoveItem(t *testing.T) {
	customerID, _ := do.NewCustomerID(1)

	tests := []struct {
		name        string
		setup       func() (*UseCase, do.OrderID, do.ItemID)
		expectedErr error
	}{
		{
			name: "ok",
			setup: func() (*UseCase, do.OrderID, do.ItemID) {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					do.Product{
						Name: "iphone",
					},
				)
				factory := do.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				orderID, _ := usecase.CreateOrder(customerID)
				itemID, _ := usecase.AddItem(orderID, "iphone")

				return usecase, orderID, itemID
			},
			expectedErr: nil,
		},
		{
			name: "remove twice",
			setup: func() (*UseCase, do.OrderID, do.ItemID) {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					do.Product{
						Name: "iphone",
					},
				)
				factory := do.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				orderID, _ := usecase.CreateOrder(customerID)
				itemID, _ := usecase.AddItem(orderID, "iphone")
				usecase.RemoveItem(orderID, itemID)

				return usecase, orderID, itemID
			},
			expectedErr: do.ErrOrderItemNotFound,
		},
		{
			name: "remove from empty",
			setup: func() (*UseCase, do.OrderID, do.ItemID) {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					do.Product{
						Name: "iphone",
					},
				)
				factory := do.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				orderID, _ := usecase.CreateOrder(customerID)

				return usecase, orderID, do.ItemID(0)
			},
			expectedErr: do.ErrOrderItemNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase, orderID, itemID := tt.setup()
			err := usecase.RemoveItem(orderID, itemID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
		})
	}
}

type FakeOrderCatalog struct {
	data map[string]do.Product
	mtx  sync.RWMutex
}

func NewFakeOrderCatalog() *FakeOrderCatalog {
	return &FakeOrderCatalog{
		data: make(map[string]do.Product),
	}
}

func (c *FakeOrderCatalog) GetProduct(name string) (do.Product, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	product, ok := c.data[name]
	if !ok {
		return do.Product{}, do.ErrProductNotFound
	}
	return product, nil
}

func (c *FakeOrderCatalog) AddProduct(product do.Product) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[product.Name] = product
}

type FakeOrderRepository struct {
	nextID do.OrderID
	orders map[do.OrderID]*do.Order
	mtx    sync.RWMutex
}

func NewFakeOrderRepository() *FakeOrderRepository {
	return &FakeOrderRepository{
		nextID: 1,
		orders: make(map[do.OrderID]*do.Order),
	}
}

func (r *FakeOrderRepository) Save(order *do.Order) error {
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

func (r *FakeOrderRepository) Update(order *do.Order) error {
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

func (r *FakeOrderRepository) NextID() (do.OrderID, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	result := r.nextID
	r.nextID += 1
	return result, nil
}

func (r *FakeOrderRepository) Get(id do.OrderID) (*do.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return nil, do.ErrOrderNotFound
	}
	return order.Clone(), nil
}

func (r *FakeOrderRepository) GetAll() ([]*do.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	result := make([]*do.Order, 0, len(r.orders))
	for _, v := range r.orders {
		result = append(result, v.Clone())
	}
	return result, nil
}
