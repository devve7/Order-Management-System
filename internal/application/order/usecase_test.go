package order

import (
	domain_order "Order-Management-System/internal/domain/order"
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

func TestAddItem(t *testing.T) {
	ctx := context.Background()
	customerID, _ := domain_order.NewCustomerID(1)

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
					domain_order.Product{
						Name: "iphone",
					},
				)
				factory := domain_order.NewOrderItemFactory(catalog)
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
					domain_order.Product{
						Name: "iphone X",
					},
				)
				factory := domain_order.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				return usecase
			},
			expectedErr: domain_order.ErrProductNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup()
			orderID, _ := usecase.CreateOrder(ctx, customerID)
			itemID, err := usecase.AddItem(ctx, orderID, "iphone")
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
			if err == nil {
				orderDTO, _ := usecase.GetOrder(ctx, orderID)
				item := orderDTO.Items[0]
				if itemID != domain_order.ItemID(item.ID) {
					t.Errorf("expected (%v), got (%v)", itemID, item.ID)
				}
			}
		})
	}
}

func TestRemoveItem(t *testing.T) {
	ctx := context.Background()
	customerID, _ := domain_order.NewCustomerID(1)

	tests := []struct {
		name        string
		setup       func() (*UseCase, domain_order.OrderID, domain_order.ItemID)
		expectedErr error
	}{
		{
			name: "ok",
			setup: func() (*UseCase, domain_order.OrderID, domain_order.ItemID) {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					domain_order.Product{
						Name: "iphone",
					},
				)
				factory := domain_order.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				orderID, _ := usecase.CreateOrder(ctx, customerID)
				itemID, _ := usecase.AddItem(ctx, orderID, "iphone")

				return usecase, orderID, itemID
			},
			expectedErr: nil,
		},
		{
			name: "remove twice",
			setup: func() (*UseCase, domain_order.OrderID, domain_order.ItemID) {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					domain_order.Product{
						Name: "iphone",
					},
				)
				factory := domain_order.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				orderID, _ := usecase.CreateOrder(ctx, customerID)
				itemID, _ := usecase.AddItem(ctx, orderID, "iphone")
				usecase.RemoveItem(ctx, orderID, itemID)

				return usecase, orderID, itemID
			},
			expectedErr: domain_order.ErrOrderItemNotFound,
		},
		{
			name: "remove from empty",
			setup: func() (*UseCase, domain_order.OrderID, domain_order.ItemID) {
				catalog := NewFakeOrderCatalog()
				catalog.AddProduct(
					domain_order.Product{
						Name: "iphone",
					},
				)
				factory := domain_order.NewOrderItemFactory(catalog)
				repo := NewFakeOrderRepository()
				usecase := NewUseCase(factory, repo)
				orderID, _ := usecase.CreateOrder(ctx, customerID)

				return usecase, orderID, domain_order.ItemID(0)
			},
			expectedErr: domain_order.ErrOrderItemNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase, orderID, itemID := tt.setup()
			err := usecase.RemoveItem(ctx, orderID, itemID)
			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected (%v), got (%v)", tt.expectedErr, err)
			}
		})
	}
}

type FakeOrderCatalog struct {
	data map[string]domain_order.Product
	mtx  sync.RWMutex
}

func NewFakeOrderCatalog() *FakeOrderCatalog {
	return &FakeOrderCatalog{
		data: make(map[string]domain_order.Product),
	}
}

func (c *FakeOrderCatalog) GetProduct(name string) (domain_order.Product, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	product, ok := c.data[name]
	if !ok {
		return domain_order.Product{}, domain_order.ErrProductNotFound
	}
	return product, nil
}

func (c *FakeOrderCatalog) AddProduct(product domain_order.Product) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[product.Name] = product
}

type FakeOrderRepository struct {
	nextID domain_order.OrderID
	orders map[domain_order.OrderID]*domain_order.Order
	mtx    sync.RWMutex
}

func NewFakeOrderRepository() *FakeOrderRepository {
	return &FakeOrderRepository{
		nextID: 1,
		orders: make(map[domain_order.OrderID]*domain_order.Order),
	}
}

func (r *FakeOrderRepository) Create(ctx context.Context, customerID domain_order.CustomerID, status domain_order.OrderStatus) (domain_order.OrderID, error) {
	orderID := r.GetNextID()
	order := domain_order.NewOrder(orderID, customerID, status, time.Now())
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.orders[orderID] = order

	return orderID, nil
}

func (r *FakeOrderRepository) Update(ctx context.Context, order *domain_order.Order) error {
	if order == nil {
		return domain_order.ErrOrderEmpty
	}
	r.mtx.Lock()
	defer r.mtx.Unlock()
	orderID := order.ID()
	if _, ok := r.orders[orderID]; !ok {
		return domain_order.ErrOrderNotFound
	}
	if order.ID() == 0 {
		return domain_order.ErrInvalidID
	}
	r.orders[orderID] = order

	return nil
}

func (r *FakeOrderRepository) GetNextID() domain_order.OrderID {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	result := r.nextID
	r.nextID += 1
	return result
}

func (r *FakeOrderRepository) Get(ctx context.Context, id domain_order.OrderID) (*domain_order.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		return nil, domain_order.ErrOrderNotFound
	}
	return order.Clone(), nil
}

func (r *FakeOrderRepository) GetAll(ctx context.Context) ([]*domain_order.Order, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	result := make([]*domain_order.Order, 0, len(r.orders))
	for _, v := range r.orders {
		result = append(result, v.Clone())
	}
	return result, nil
}
