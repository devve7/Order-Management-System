package product

import (
	"Order-Management-System/internal/application/ports"
	domain_product "Order-Management-System/internal/domain/product"
	"context"
	"errors"
	"testing"
)

type repoMock struct {
	createFn func(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error)
	getFn    func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error)
	getAllFn func(ctx context.Context) ([]*domain_product.Product, error)
	updateFn func(ctx context.Context, product *domain_product.Product) error
}

func (m *repoMock) Create(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
	if m.createFn != nil {
		return m.createFn(ctx, product)
	}
	return 0, nil
}

func (m *repoMock) Get(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, nil
}

func (m *repoMock) List(ctx context.Context, params domain_product.ProductListParams) ([]*domain_product.Product, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}

func (m *repoMock) Update(ctx context.Context, product *domain_product.Product) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, product)
	}
	return nil
}

func mustProductID(t *testing.T, id int64) domain_product.ProductID {
	t.Helper()
	v, err := domain_product.NewProductID(id)
	if err != nil {
		t.Fatalf("NewProductID(%d) error = %v", id, err)
	}
	return v
}

func mustProductName(t *testing.T, name string) domain_product.ProductName {
	t.Helper()
	v, err := domain_product.NewProductName(name)
	if err != nil {
		t.Fatalf("NewProductName(%q) error = %v", name, err)
	}
	return v
}

func mustPrice(t *testing.T, price int64) domain_product.Price {
	t.Helper()
	v, err := domain_product.NewPrice(price)
	if err != nil {
		t.Fatalf("NewPrice(%d) error = %v", price, err)
	}
	return v
}

func mustStock(t *testing.T, stock int64) domain_product.Stock {
	t.Helper()
	v, err := domain_product.NewStock(stock)
	if err != nil {
		t.Fatalf("NewStock(%d) error = %v", stock, err)
	}
	return v
}

func TestNewUseCase(t *testing.T) {
	cache := newMockCache()
	repo := &repoMock{}
	u := NewUseCase(repo, cache)

	if u == nil {
		t.Fatal("NewUseCase() returned nil")
	}
	if u.repo != repo {
		t.Fatal("repo was not assigned")
	}
}

func TestCreateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		called := false

		repo := &repoMock{
			createFn: func(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
				called = true

				if product.ID() != 0 {
					t.Fatalf("product.ID() = %v, want 0", product.ID())
				}
				if product.Name() != "iPhone" {
					t.Fatalf("product.Name() = %v, want %v", product.Name(), domain_product.ProductName("iPhone"))
				}
				if product.Price() != 100000 {
					t.Fatalf("product.Price() = %v, want %v", product.Price(), domain_product.Price(100000))
				}
				if product.Stock() != 10 {
					t.Fatalf("product.Stock() = %v, want %v", product.Stock(), domain_product.Stock(10))
				}
				if !product.Active() {
					t.Fatal("product.Active() = false, want true")
				}

				return 42, nil
			},
		}

		cache := newMockCache()
		u := NewUseCase(repo, cache)

		got, err := u.CreateProduct(context.Background(), "iPhone", 100000, 10)
		if err != nil {
			t.Fatalf("CreateProduct() error = %v", err)
		}
		if !called {
			t.Fatal("repo.Create was not called")
		}
		if got != 42 {
			t.Fatalf("CreateProduct() = %d, want 42", got)
		}
	})

	t.Run("invalid name", func(t *testing.T) {
		repoCalled := false
		repo := &repoMock{
			createFn: func(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
				repoCalled = true
				return 0, nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		_, err := u.CreateProduct(context.Background(), "", 100, 1)
		if err != domain_product.ErrInvalidProductName {
			t.Fatalf("CreateProduct() error = %v, want %v", err, domain_product.ErrInvalidProductName)
		}
		if repoCalled {
			t.Fatal("repo.Create should not be called")
		}
	})

	t.Run("invalid price", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		_, err := u.CreateProduct(context.Background(), "iPhone", -1, 1)
		if err != domain_product.ErrInvalidPrice {
			t.Fatalf("CreateProduct() error = %v, want %v", err, domain_product.ErrInvalidPrice)
		}
	})

	t.Run("invalid stock", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		_, err := u.CreateProduct(context.Background(), "iPhone", 100, -1)
		if err != domain_product.ErrInvalidStock {
			t.Fatalf("CreateProduct() error = %v, want %v", err, domain_product.ErrInvalidStock)
		}
	})

	t.Run("repo create error", func(t *testing.T) {
		wantErr := errors.New("repo create failed")
		repo := &repoMock{
			createFn: func(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
				return 0, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		_, err := u.CreateProduct(context.Background(), "iPhone", 100, 1)
		if err != wantErr {
			t.Fatalf("CreateProduct() error = %v, want %v", err, wantErr)
		}
	})
}

func TestGetProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 7),
			mustProductName(t, "MacBook"),
			mustPrice(t, 200000),
			mustStock(t, 5),
			false,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				if id != 7 {
					t.Fatalf("id = %v, want %v", id, domain_product.ProductID(7))
				}
				return product, nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		dto, err := u.GetProduct(context.Background(), 7)
		if err != nil {
			t.Fatalf("GetProduct() error = %v", err)
		}
		if dto.ID != 7 || dto.Name != "MacBook" || dto.Price != 200000 || dto.Stock != 5 || dto.Active != false {
			t.Fatalf("unexpected dto: %+v", dto)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		_, err := u.GetProduct(context.Background(), 0)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("GetProduct() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		_, err := u.GetProduct(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("GetProduct() error = %v, want %v", err, wantErr)
		}
	})
}

func TestChangePrice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		updateCalled := false
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				updateCalled = true
				if updated.Price() != 120000 {
					t.Fatalf("updated.Price() = %v, want %v", updated.Price(), domain_product.Price(120000))
				}
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.ChangePrice(context.Background(), 1, 120000)
		if err != nil {
			t.Fatalf("ChangePrice() error = %v", err)
		}
		if !updateCalled {
			t.Fatal("repo.Update was not called")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.ChangePrice(context.Background(), 0, 100)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("ChangePrice() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("invalid price", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.ChangePrice(context.Background(), 1, -1)
		if err != domain_product.ErrInvalidPrice {
			t.Fatalf("ChangePrice() error = %v, want %v", err, domain_product.ErrInvalidPrice)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.ChangePrice(context.Background(), 1, 200)
		if err != wantErr {
			t.Fatalf("ChangePrice() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		wantErr := errors.New("repo update failed")
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				return wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.ChangePrice(context.Background(), 1, 120000)
		if err != wantErr {
			t.Fatalf("ChangePrice() error = %v, want %v", err, wantErr)
		}
	})
}

func TestAddStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				if updated.Stock() != 15 {
					t.Fatalf("updated.Stock() = %v, want %v", updated.Stock(), domain_product.Stock(15))
				}
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.AddStock(context.Background(), 1, 5)
		if err != nil {
			t.Fatalf("AddStock() error = %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.AddStock(context.Background(), 0, 5)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("AddStock() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("invalid stock", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.AddStock(context.Background(), 1, -1)
		if err != domain_product.ErrInvalidStock {
			t.Fatalf("AddStock() error = %v, want %v", err, domain_product.ErrInvalidStock)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.AddStock(context.Background(), 1, 5)
		if err != wantErr {
			t.Fatalf("AddStock() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		wantErr := errors.New("repo update failed")
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				return wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.AddStock(context.Background(), 1, 5)
		if err != wantErr {
			t.Fatalf("AddStock() error = %v, want %v", err, wantErr)
		}
	})
}

func TestRemoveStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		updateCalled := false
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				updateCalled = true
				if updated.Stock() != 6 {
					t.Fatalf("updated.Stock() = %v, want %v", updated.Stock(), domain_product.Stock(6))
				}
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.RemoveStock(context.Background(), 1, 4)
		if err != nil {
			t.Fatalf("RemoveStock() error = %v", err)
		}
		if !updateCalled {
			t.Fatal("repo.Update was not called")
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.RemoveStock(context.Background(), 0, 1)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("RemoveStock() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("invalid stock", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.RemoveStock(context.Background(), 1, -1)
		if err != domain_product.ErrInvalidStock {
			t.Fatalf("RemoveStock() error = %v, want %v", err, domain_product.ErrInvalidStock)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.RemoveStock(context.Background(), 1, 4)
		if err != wantErr {
			t.Fatalf("RemoveStock() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("domain remove stock error", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				t.Fatal("repo.Update should not be called")
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.RemoveStock(context.Background(), 1, 11)
		if err != domain_product.ErrInsufficientStock {
			t.Fatalf("RemoveStock() error = %v, want %v", err, domain_product.ErrInsufficientStock)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		wantErr := errors.New("repo update failed")
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				return wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.RemoveStock(context.Background(), 1, 4)
		if err != wantErr {
			t.Fatalf("RemoveStock() error = %v, want %v", err, wantErr)
		}
	})
}

func TestDeactivateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				if updated.Active() {
					t.Fatal("updated.Active() = true, want false")
				}
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.DeactivateProduct(context.Background(), 1)
		if err != nil {
			t.Fatalf("DeactivateProduct() error = %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.DeactivateProduct(context.Background(), 0)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("DeactivateProduct() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.DeactivateProduct(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("DeactivateProduct() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("domain deactivate error", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				t.Fatal("repo.Update should not be called")
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.DeactivateProduct(context.Background(), 1)
		if err != domain_product.ErrInactiveProduct {
			t.Fatalf("DeactivateProduct() error = %v, want %v", err, domain_product.ErrInactiveProduct)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		wantErr := errors.New("repo update failed")
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				return wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.DeactivateProduct(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("DeactivateProduct() error = %v, want %v", err, wantErr)
		}
	})
}

func TestActivateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				if !updated.Active() {
					t.Fatal("updated.Active() = false, want true")
				}
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.ActivateProduct(context.Background(), 1)
		if err != nil {
			t.Fatalf("ActivateProduct() error = %v", err)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.ActivateProduct(context.Background(), 0)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("ActivateProduct() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.ActivateProduct(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("ActivateProduct() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("domain activate error", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				t.Fatal("repo.Update should not be called")
				return nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.ActivateProduct(context.Background(), 1)
		if err != domain_product.ErrProductAlreadyActive {
			t.Fatalf("ActivateProduct() error = %v, want %v", err, domain_product.ErrProductAlreadyActive)
		}
	})

	t.Run("repo update error", func(t *testing.T) {
		wantErr := errors.New("repo update failed")
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				return wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.ActivateProduct(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("ActivateProduct() error = %v, want %v", err, wantErr)
		}
	})
}

func TestGetSnapshot(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 7),
			mustProductName(t, "MacBook"),
			mustPrice(t, 200000),
			mustStock(t, 5),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		got, err := u.GetSnapshot(context.Background(), 7)
		if err != nil {
			t.Fatalf("GetSnapshot() error = %v", err)
		}

		want := ports.ProductSnapshot{
			ProductID: 7,
			Name:      "MacBook",
			Price:     200000,
		}

		if got != want {
			t.Fatalf("GetSnapshot() = %+v, want %+v", got, want)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		_, err := u.GetSnapshot(context.Background(), 0)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("GetSnapshot() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		_, err := u.GetSnapshot(context.Background(), 1)
		if err != wantErr {
			t.Fatalf("GetSnapshot() error = %v, want %v", err, wantErr)
		}
	})
}

func TestEnsureAvailable(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.EnsureAvailable(context.Background(), 1, 5)
		if err != nil {
			t.Fatalf("EnsureAvailable() error = %v", err)
		}
	})

	t.Run("invalid product id", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.EnsureAvailable(context.Background(), 0, 5)
		if err != domain_product.ErrInvalidProductID {
			t.Fatalf("EnsureAvailable() error = %v, want %v", err, domain_product.ErrInvalidProductID)
		}
	})

	t.Run("invalid quantity", func(t *testing.T) {
		cache := newMockCache()
		u := NewUseCase(&repoMock{}, cache)

		err := u.EnsureAvailable(context.Background(), 1, -1)
		if err != domain_product.ErrInvalidStock {
			t.Fatalf("EnsureAvailable() error = %v, want %v", err, domain_product.ErrInvalidStock)
		}
	})

	t.Run("repo get error", func(t *testing.T) {
		wantErr := errors.New("repo get failed")
		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, wantErr
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.EnsureAvailable(context.Background(), 1, 5)
		if err != wantErr {
			t.Fatalf("EnsureAvailable() error = %v, want %v", err, wantErr)
		}
	})

	t.Run("inactive product", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.EnsureAvailable(context.Background(), 1, 5)
		if err != domain_product.ErrInactiveProduct {
			t.Fatalf("EnsureAvailable() error = %v, want %v", err, domain_product.ErrInactiveProduct)
		}
	})

	t.Run("insufficient stock", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &repoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
		}
		cache := newMockCache()
		u := NewUseCase(repo, cache)

		err := u.EnsureAvailable(context.Background(), 1, 11)
		if err != domain_product.ErrInsufficientStock {
			t.Fatalf("EnsureAvailable() error = %v, want %v", err, domain_product.ErrInsufficientStock)
		}
	})
}
