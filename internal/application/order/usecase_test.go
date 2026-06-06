package order

import (
	"Order-Management-System/internal/application/ports"
	domain_order "Order-Management-System/internal/domain/order"
	"context"
	"errors"
	"testing"
	"time"
)

type repoMock struct {
	createFn func(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error)
	getFn    func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error)
	getAllFn func(ctx context.Context) ([]*domain_order.Order, error)
	updateFn func(ctx context.Context, order *domain_order.Order) error
}

func (m *repoMock) Create(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
	if m.createFn != nil {
		return m.createFn(ctx, customerID)
	}
	return 0, nil
}

func (m *repoMock) Get(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
	if m.getFn != nil {
		return m.getFn(ctx, orderID)
	}
	return nil, nil
}

func (m *repoMock) List(ctx context.Context, params domain_order.OrderListParams) ([]*domain_order.Order, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}

func (m *repoMock) Update(ctx context.Context, order *domain_order.Order) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, order)
	}
	return nil
}

// Здесь оставляем пустой стаб, потому что для большинства тестов productService не нужен.
// Для AddItem нужен точный интерфейс ports.ProductForOrder, а его сигнатура не загружена.
type productServiceStub struct{}

func mustOrderID(t *testing.T, id int64) domain_order.OrderID {
	t.Helper()
	v, err := domain_order.NewOrderID(id)
	if err != nil {
		t.Fatalf("NewOrderID(%d) error = %v", id, err)
	}
	return v
}

func mustCustomerID(t *testing.T, id int64) domain_order.CustomerID {
	t.Helper()
	v, err := domain_order.NewCustomerID(id)
	if err != nil {
		t.Fatalf("NewCustomerID(%d) error = %v", id, err)
	}
	return v
}

func mustProductID(t *testing.T, id int64) domain_order.ProductID {
	t.Helper()
	v, err := domain_order.NewProductID(id)
	if err != nil {
		t.Fatalf("NewProductID(%d) error = %v", id, err)
	}
	return v
}

func mustProductName(t *testing.T, s string) domain_order.ProductName {
	t.Helper()
	v, err := domain_order.NewProductName(s)
	if err != nil {
		t.Fatalf("NewProductName(%q) error = %v", s, err)
	}
	return v
}

func mustPrice(t *testing.T, amount int64) domain_order.Price {
	t.Helper()
	v, err := domain_order.NewPrice(amount)
	if err != nil {
		t.Fatalf("NewPrice(%d) error = %v", amount, err)
	}
	return v
}

func mustQuantity(t *testing.T, q int64) domain_order.Quantity {
	t.Helper()
	v, err := domain_order.NewQuantity(q)
	if err != nil {
		t.Fatalf("NewQuantity(%d) error = %v", q, err)
	}
	return v
}

func mustItemID(t *testing.T, id int64) domain_order.ItemID {
	t.Helper()
	v, err := domain_order.NewItemID(id)
	if err != nil {
		t.Fatalf("NewItemID(%d) error = %v", id, err)
	}
	return v
}

func newTestOrder(t *testing.T) *domain_order.Order {
	t.Helper()
	return domain_order.NewOrder(
		mustOrderID(t, 1),
		mustCustomerID(t, 10),
		domain_order.StatusCreated,
		time.Date(2026, 4, 22, 12, 0, 0, 0, time.UTC),
	)
}

func addTestItem(t *testing.T, order *domain_order.Order) {
	t.Helper()
	err := order.AddItem(
		mustProductID(t, 100),
		mustProductName(t, "apple"),
		mustPrice(t, 150),
		mustQuantity(t, 2),
	)
	if err != nil {
		t.Fatalf("order.AddItem() error = %v", err)
	}
}

func TestNewUseCase(t *testing.T) {
	repo := &repoMock{}
	u := NewUseCase(nil, repo)

	if u == nil {
		t.Fatal("NewUseCase() returned nil")
	}
	if u.repo != repo {
		t.Fatal("repo was not assigned")
	}
	if u.productService != nil {
		t.Fatal("productService should be nil")
	}
}

func TestCreateOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		called := false
		repo := &repoMock{
			createFn: func(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
				called = true
				if customerID != 42 {
					t.Fatalf("customerID = %v, want %v", customerID, domain_order.CustomerID(42))
				}
				return 100, nil
			},
		}

		u := NewUseCase(nil, repo)

		got, err := u.CreateOrder(context.Background(), 42)
		if err != nil {
			t.Fatalf("CreateOrder() error = %v", err)
		}
		if !called {
			t.Fatal("repo.Create was not called")
		}
		if got != 100 {
			t.Fatalf("CreateOrder() = %v, want %v", got, domain_order.OrderID(100))
		}
	})

	t.Run("invalid customer id", func(t *testing.T) {
		repoCalled := false
		repo := &repoMock{
			createFn: func(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
				repoCalled = true
				return 0, nil
			},
		}

		u := NewUseCase(nil, repo)

		_, err := u.CreateOrder(context.Background(), -1)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("CreateOrder() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
		if repoCalled {
			t.Fatal("repo.Create should not be called")
		}
	})

	t.Run("repo create error", func(t *testing.T) {
		wantErr := errors.New("repo create failed")
		repo := &repoMock{
			createFn: func(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
				return 0, wantErr
			},
		}

		u := NewUseCase(nil, repo)

		_, err := u.CreateOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("CreateOrder() error = %v, want %v", err, wantErr)
		}
	})
}

func TestRemoveItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		getCalled := false
		updateCalled := false

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				getCalled = true
				if orderID != 1 {
					t.Fatalf("orderID = %v, want %v", orderID, domain_order.OrderID(1))
				}
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				updateCalled = true
				if len(updated.Items()) != 0 {
					t.Fatalf("len(updated.Items()) = %d, want 0", len(updated.Items()))
				}
				return nil
			},
		}

		u := NewUseCase(nil, repo)

		err := u.RemoveItem(context.Background(), 1, 1)
		if err != nil {
			t.Fatalf("RemoveItem() error = %v", err)
		}
		if !getCalled {
			t.Fatal("repo.Get was not called")
		}
		if !updateCalled {
			t.Fatal("repo.Update was not called")
		}
	})

	t.Run("invalid order id", func(t *testing.T) {
		repo := &repoMock{}
		u := NewUseCase(nil, repo)

		err := u.RemoveItem(context.Background(), -1, 1)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("RemoveItem() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
	})

	t.Run("invalid item id", func(t *testing.T) {
		repo := &repoMock{}
		u := NewUseCase(nil, repo)

		err := u.RemoveItem(context.Background(), 1, -1)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("RemoveItem() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return nil, wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.RemoveItem(context.Background(), 1, 1)
		if err != wantErr {
			t.Fatalf("RemoveItem() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("domain remove error", func(t *testing.T) {
		order := newTestOrder(t)

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}
		u := NewUseCase(nil, repo)

		err := u.RemoveItem(context.Background(), 1, 1)
		if err != domain_order.ErrOrderItemNotFound {
			t.Fatalf("RemoveItem() error = %v, want %v", err, domain_order.ErrOrderItemNotFound)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		wantErr := errors.New("repo update failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, order *domain_order.Order) error {
				return wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.RemoveItem(context.Background(), 1, 1)
		if err != wantErr {
			t.Fatalf("RemoveItem() error = %v, want %v", err, wantErr)
		}
	})
}

func TestPayOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		updateCalled := false
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				updateCalled = true
				if updated.Status() != domain_order.StatusPaid {
					t.Fatalf("updated.Status() = %v, want %v", updated.Status(), domain_order.StatusPaid)
				}
				return nil
			},
		}

		u := NewUseCase(nil, repo)

		err := u.PayOrder(context.Background(), 1)
		if err != nil {
			t.Fatalf("PayOrder() error = %v", err)
		}
		if !updateCalled {
			t.Fatal("repo.Update was not called")
		}
	})

	t.Run("invalid order id", func(t *testing.T) {
		u := NewUseCase(nil, &repoMock{})

		err := u.PayOrder(context.Background(), -1)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("PayOrder() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return nil, wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.PayOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("PayOrder() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("domain pay error", func(t *testing.T) {
		order := newTestOrder(t)

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}
		u := NewUseCase(nil, repo)

		err := u.PayOrder(context.Background(), 1)
		if err != domain_order.ErrOrderEmpty {
			t.Fatalf("PayOrder() error = %v, want %v", err, domain_order.ErrOrderEmpty)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		wantErr := errors.New("repo update failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, order *domain_order.Order) error {
				return wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.PayOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("PayOrder() error = %v, want %v", err, wantErr)
		}
	})
}

func TestShipOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := domain_order.NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 10),
			domain_order.StatusPaid,
			time.Now(),
		)

		updateCalled := false
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				updateCalled = true
				if updated.Status() != domain_order.StatusShipped {
					t.Fatalf("updated.Status() = %v, want %v", updated.Status(), domain_order.StatusShipped)
				}
				return nil
			},
		}

		u := NewUseCase(nil, repo)

		err := u.ShipOrder(context.Background(), 1)
		if err != nil {
			t.Fatalf("ShipOrder() error = %v", err)
		}
		if !updateCalled {
			t.Fatal("repo.Update was not called")
		}
	})

	t.Run("invalid order id", func(t *testing.T) {
		u := NewUseCase(nil, &repoMock{})

		err := u.ShipOrder(context.Background(), -1)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("ShipOrder() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return nil, wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.ShipOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("ShipOrder() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("domain ship error", func(t *testing.T) {
		order := newTestOrder(t)

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}
		u := NewUseCase(nil, repo)

		err := u.ShipOrder(context.Background(), 1)
		if err != domain_order.ErrCannotShip {
			t.Fatalf("ShipOrder() error = %v, want %v", err, domain_order.ErrCannotShip)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		order := domain_order.NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 10),
			domain_order.StatusPaid,
			time.Now(),
		)

		wantErr := errors.New("repo update failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, order *domain_order.Order) error {
				return wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.ShipOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("ShipOrder() error = %v, want %v", err, wantErr)
		}
	})
}

func TestCancelOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)

		updateCalled := false
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				updateCalled = true
				if updated.Status() != domain_order.StatusCancelled {
					t.Fatalf("updated.Status() = %v, want %v", updated.Status(), domain_order.StatusCancelled)
				}
				return nil
			},
		}

		u := NewUseCase(nil, repo)

		err := u.CancelOrder(context.Background(), 1)
		if err != nil {
			t.Fatalf("CancelOrder() error = %v", err)
		}
		if !updateCalled {
			t.Fatal("repo.Update was not called")
		}
	})

	t.Run("invalid order id", func(t *testing.T) {
		u := NewUseCase(nil, &repoMock{})

		err := u.CancelOrder(context.Background(), -1)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("CancelOrder() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return nil, wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.CancelOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("CancelOrder() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("domain cancel error", func(t *testing.T) {
		order := domain_order.NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 10),
			domain_order.StatusShipped,
			time.Now(),
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}
		u := NewUseCase(nil, repo)

		err := u.CancelOrder(context.Background(), 1)
		if err != domain_order.ErrCannotCancel {
			t.Fatalf("CancelOrder() error = %v, want %v", err, domain_order.ErrCannotCancel)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		order := newTestOrder(t)

		wantErr := errors.New("repo update failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, order *domain_order.Order) error {
				return wantErr
			},
		}
		u := NewUseCase(nil, repo)

		err := u.CancelOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("CancelOrder() error = %v, want %v", err, wantErr)
		}
	})
}

func TestGetOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}
		u := NewUseCase(nil, repo)

		got, err := u.GetOrder(context.Background(), 1)
		if err != nil {
			t.Fatalf("GetOrder() error = %v", err)
		}

		if got.ID != 1 {
			t.Fatalf("dto.ID = %d, want 1", got.ID)
		}
		if got.CustomerID != 10 {
			t.Fatalf("dto.CustomerID = %d, want 10", got.CustomerID)
		}
		if got.Status != "created" {
			t.Fatalf("dto.Status = %q, want %q", got.Status, "created")
		}
		if len(got.Items) != 1 {
			t.Fatalf("len(dto.Items) = %d, want 1", len(got.Items))
		}
		if got.Total != 300 {
			t.Fatalf("dto.Total = %d, want 300", got.Total)
		}
	})

	t.Run("invalid order id", func(t *testing.T) {
		u := NewUseCase(nil, &repoMock{})

		_, err := u.GetOrder(context.Background(), -1)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("GetOrder() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return nil, wantErr
			},
		}
		u := NewUseCase(nil, repo)

		_, err := u.GetOrder(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("GetOrder() error = %v, want %v", err, wantErr)
		}
	})
}

type productServiceMock struct {
	ensureAvailableFn func(ctx context.Context, productID int64, quantity int64) error
	getSnapshotFn     func(ctx context.Context, productID int64) (ports.ProductSnapshot, error)
}

func (m *productServiceMock) EnsureAvailable(ctx context.Context, productID int64, quantity int64) error {
	if m.ensureAvailableFn != nil {
		return m.ensureAvailableFn(ctx, productID, quantity)
	}
	return nil
}

func (m *productServiceMock) GetSnapshot(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
	if m.getSnapshotFn != nil {
		return m.getSnapshotFn(ctx, productID)
	}
	return ports.ProductSnapshot{}, nil
}

func TestAddItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)

		ensureCalled := false
		snapshotCalled := false
		getCalled := false
		updateCalled := false

		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				ensureCalled = true
				if productID != 100 {
					t.Fatalf("productID = %d, want 100", productID)
				}
				if quantity != 2 {
					t.Fatalf("quantity = %d, want 2", quantity)
				}
				return nil
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				snapshotCalled = true
				if productID != 100 {
					t.Fatalf("productID = %d, want 100", productID)
				}
				return ports.ProductSnapshot{
					ProductID: 100,
					Name:      "apple",
					Price:     150,
				}, nil
			},
		}

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				getCalled = true
				if orderID != 1 {
					t.Fatalf("orderID = %v, want %v", orderID, domain_order.OrderID(1))
				}
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				updateCalled = true

				items := updated.Items()
				if len(items) != 1 {
					t.Fatalf("len(updated.Items()) = %d, want 1", len(items))
				}
				if items[0].ProductID() != 100 {
					t.Fatalf("items[0].ProductID() = %v, want %v", items[0].ProductID(), domain_order.ProductID(100))
				}
				if items[0].Name() != "apple" {
					t.Fatalf("items[0].Name() = %v, want %v", items[0].Name(), domain_order.ProductName("apple"))
				}
				if items[0].Price() != 150 {
					t.Fatalf("items[0].Price() = %v, want %v", items[0].Price(), domain_order.Price(150))
				}
				if items[0].Quantity() != 2 {
					t.Fatalf("items[0].Quantity() = %v, want %v", items[0].Quantity(), domain_order.Quantity(2))
				}
				return nil
			},
		}

		u := NewUseCase(productService, repo)

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != nil {
			t.Fatalf("AddItem() error = %v", err)
		}
		if !ensureCalled {
			t.Fatal("productService.EnsureAvailable was not called")
		}
		if !snapshotCalled {
			t.Fatal("productService.GetSnapshot was not called")
		}
		if !getCalled {
			t.Fatal("repo.Get was not called")
		}
		if !updateCalled {
			t.Fatal("repo.Update was not called")
		}
	})

	t.Run("invalid order id", func(t *testing.T) {
		productServiceCalled := false
		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				productServiceCalled = true
				return nil
			},
		}

		u := NewUseCase(productService, &repoMock{})

		err := u.AddItem(context.Background(), -1, 100, 2)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("AddItem() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
		if productServiceCalled {
			t.Fatal("productService should not be called")
		}
	})

	t.Run("invalid product id", func(t *testing.T) {
		productServiceCalled := false
		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				productServiceCalled = true
				return nil
			},
		}

		u := NewUseCase(productService, &repoMock{})

		err := u.AddItem(context.Background(), 1, -1, 2)
		if err != domain_order.ErrInvalidID {
			t.Fatalf("AddItem() error = %v, want %v", err, domain_order.ErrInvalidID)
		}
		if productServiceCalled {
			t.Fatal("productService should not be called")
		}
	})

	t.Run("invalid quantity", func(t *testing.T) {
		productServiceCalled := false
		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				productServiceCalled = true
				return nil
			},
		}

		u := NewUseCase(productService, &repoMock{})

		err := u.AddItem(context.Background(), 1, 100, 0)
		if err != domain_order.ErrInvalidQuantity {
			t.Fatalf("AddItem() error = %v, want %v", err, domain_order.ErrInvalidQuantity)
		}
		if productServiceCalled {
			t.Fatal("productService should not be called")
		}
	})

	t.Run("ensure available error", func(t *testing.T) {
		wantErr := errors.New("ensure available failed")

		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return wantErr
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				t.Fatal("GetSnapshot should not be called")
				return ports.ProductSnapshot{}, nil
			},
		}

		u := NewUseCase(productService, &repoMock{})

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != wantErr {
			t.Fatalf("AddItem() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("get snapshot error", func(t *testing.T) {
		wantErr := errors.New("get snapshot failed")

		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return nil
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				return ports.ProductSnapshot{}, wantErr
			},
		}

		u := NewUseCase(productService, &repoMock{})

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != wantErr {
			t.Fatalf("AddItem() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")

		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return nil
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				return ports.ProductSnapshot{
					ProductID: 100,
					Name:      "apple",
					Price:     150,
				}, nil
			},
		}

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return nil, wantErr
			},
		}

		u := NewUseCase(productService, repo)

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != wantErr {
			t.Fatalf("AddItem() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("invalid product name in snapshot", func(t *testing.T) {
		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return nil
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				return ports.ProductSnapshot{
					ProductID: 100,
					Name:      "",
					Price:     150,
				}, nil
			},
		}

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return newTestOrder(t), nil
			},
		}

		u := NewUseCase(productService, repo)

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != domain_order.ErrInvalidProductName {
			t.Fatalf("AddItem() error = %v, want %v", err, domain_order.ErrInvalidProductName)
		}
	})

	t.Run("invalid product price in snapshot", func(t *testing.T) {
		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return nil
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				return ports.ProductSnapshot{
					ProductID: 100,
					Name:      "apple",
					Price:     -1,
				}, nil
			},
		}

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return newTestOrder(t), nil
			},
		}

		u := NewUseCase(productService, repo)

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != domain_order.ErrInvalidPrice {
			t.Fatalf("AddItem() error = %v, want %v", err, domain_order.ErrInvalidPrice)
		}
	})

	t.Run("domain add item error for paid order", func(t *testing.T) {
		order := domain_order.NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 10),
			domain_order.StatusPaid,
			time.Now(),
		)

		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return nil
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				return ports.ProductSnapshot{
					ProductID: 100,
					Name:      "apple",
					Price:     150,
				}, nil
			},
		}

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, order *domain_order.Order) error {
				t.Fatal("repo.Update should not be called")
				return nil
			},
		}

		u := NewUseCase(productService, repo)

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != domain_order.ErrOrderPaid {
			t.Fatalf("AddItem() error = %v, want %v", err, domain_order.ErrOrderPaid)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		order := newTestOrder(t)
		wantErr := errors.New("repo update failed")

		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return nil
			},
			getSnapshotFn: func(ctx context.Context, productID int64) (ports.ProductSnapshot, error) {
				return ports.ProductSnapshot{
					ProductID: 100,
					Name:      "apple",
					Price:     150,
				}, nil
			},
		}

		repo := &repoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				return wantErr
			},
		}

		u := NewUseCase(productService, repo)

		err := u.AddItem(context.Background(), 1, 100, 2)
		if err != wantErr {
			t.Fatalf("AddItem() error = %v, want %v", err, wantErr)
		}
	})
}
