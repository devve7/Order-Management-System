package order

import (
	application_order "Order-Management-System/internal/application/order"
	"Order-Management-System/internal/application/ports"
	domain_order "Order-Management-System/internal/domain/order"
	domain_product "Order-Management-System/internal/domain/product"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type orderRepoMock struct {
	createFn func(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error)
	getFn    func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error)
	getAllFn func(ctx context.Context) ([]*domain_order.Order, error)
	updateFn func(ctx context.Context, order *domain_order.Order) error
}

func (m *orderRepoMock) Create(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
	if m.createFn != nil {
		return m.createFn(ctx, customerID)
	}
	return 0, nil
}

func (m *orderRepoMock) Get(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
	if m.getFn != nil {
		return m.getFn(ctx, orderID)
	}
	return nil, nil
}

func (m *orderRepoMock) List(ctx context.Context, params domain_order.OrderListParams) ([]*domain_order.Order, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}

func (m *orderRepoMock) Update(ctx context.Context, order *domain_order.Order) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, order)
	}
	return nil
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

func newOrderHandler(productService ports.ProductForOrder, repo domain_order.Repository) *OrderHandler {
	usecase := application_order.NewUseCase(productService, repo)
	logger := logrus.New()
	logger.SetOutput(bytes.NewBuffer(nil))
	return NewOrderHandler(usecase, logger)
}

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

func TestOrderHandler_CreateOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &orderRepoMock{
			createFn: func(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
				if customerID != 42 {
					t.Fatalf("customerID = %v, want %v", customerID, domain_order.CustomerID(42))
				}
				return 100, nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{"customer_id":42}`))
		w := httptest.NewRecorder()

		h.CreateOrder(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
		}
		if ct := w.Header().Get("Content-Type"); ct != "application/json" {
			t.Fatalf("Content-Type = %q, want %q", ct, "application/json")
		}

		var got OrderIDDTO
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if got.OrderID != 100 {
			t.Fatalf("order_id = %d, want 100", got.OrderID)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{"customer_id":`))
		w := httptest.NewRecorder()

		h.CreateOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing customer id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{}`))
		w := httptest.NewRecorder()

		h.CreateOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var got ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if got.Error != "customer_id is required" {
			t.Fatalf("error = %q, want %q", got.Error, "customer_id is required")
		}
	})

	t.Run("unexpected repo error becomes 500", func(t *testing.T) {
		repo := &orderRepoMock{
			createFn: func(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
				return 0, errors.New("db exploded")
			},
		}
		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{"customer_id":42}`))
		w := httptest.NewRecorder()

		h.CreateOrder(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		var got ErrorResponse
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("decode error response: %v", err)
		}
		if got.Error != "internal server error" {
			t.Fatalf("error = %q, want %q", got.Error, "internal server error")
		}
	})
}

func TestOrderHandler_GetOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				if orderID != 1 {
					t.Fatalf("orderID = %v, want %v", orderID, domain_order.OrderID(1))
				}
				return order, nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetOrder(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var got OrderDTO
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if got.ID != 1 {
			t.Fatalf("order_id = %d, want 1", got.ID)
		}
		if got.CustomerID != 10 {
			t.Fatalf("customer_id = %d, want 10", got.CustomerID)
		}
		if got.Status != "created" {
			t.Fatalf("status = %q, want %q", got.Status, "created")
		}
		if len(got.Items) != 1 {
			t.Fatalf("len(items) = %d, want 1", len(got.Items))
		}
		if got.Total != 300 {
			t.Fatalf("total = %d, want 300", got.Total)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodGet, "/orders/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return nil, domain_order.ErrOrderNotFound
			},
		}
		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetOrder(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestOrderHandler_GetOrders(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order1 := newTestOrder(t)
		addTestItem(t, order1)

		order2 := domain_order.NewOrder(
			mustOrderID(t, 2),
			mustCustomerID(t, 20),
			domain_order.StatusPaid,
			time.Date(2026, 4, 22, 13, 0, 0, 0, time.UTC),
		)

		repo := &orderRepoMock{
			getAllFn: func(ctx context.Context) ([]*domain_order.Order, error) {
				return []*domain_order.Order{order1, order2}, nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()

		h.GetOrders(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var got []OrderDTO
		if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(got) != 2 {
			t.Fatalf("len(got) = %d, want 2", len(got))
		}
		if got[0].ID != 1 || got[1].ID != 2 {
			t.Fatalf("ids = [%d %d], want [1 2]", got[0].ID, got[1].ID)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &orderRepoMock{
			getAllFn: func(ctx context.Context) ([]*domain_order.Order, error) {
				return nil, errors.New("db failed")
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()

		h.GetOrders(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestOrderHandler_AddItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)

		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				if productID != 100 {
					t.Fatalf("productID = %d, want 100", productID)
				}
				if quantity != 2 {
					t.Fatalf("quantity = %d, want 2", quantity)
				}
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

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				if len(updated.Items()) != 1 {
					t.Fatalf("len(updated.Items()) = %d, want 1", len(updated.Items()))
				}
				return nil
			},
		}

		h := newOrderHandler(productService, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/items", bytes.NewBufferString(`{"product_id":100,"quantity":2}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddItem(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/1/items", bytes.NewBufferString(`{"product_id":`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddItem(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing product id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/1/items", bytes.NewBufferString(`{"quantity":2}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddItem(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing quantity", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/1/items", bytes.NewBufferString(`{"product_id":100}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddItem(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("bad order id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/abc/items", bytes.NewBufferString(`{"product_id":100,"quantity":2}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.AddItem(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("conflict from product service", func(t *testing.T) {
		productService := &productServiceMock{
			ensureAvailableFn: func(ctx context.Context, productID int64, quantity int64) error {
				return domain_product.ErrInsufficientStock
			},
		}
		h := newOrderHandler(productService, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/1/items", bytes.NewBufferString(`{"product_id":100,"quantity":2}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddItem(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})

	t.Run("paid order returns conflict", func(t *testing.T) {
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

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}

		h := newOrderHandler(productService, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/items", bytes.NewBufferString(`{"product_id":100,"quantity":2}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddItem(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})
}

func TestOrderHandler_DeleteItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				if len(updated.Items()) != 0 {
					t.Fatalf("len(updated.Items()) = %d, want 0", len(updated.Items()))
				}
				return nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodDelete, "/orders/1/items/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1", "item_id": "1"})
		w := httptest.NewRecorder()

		h.DeleteItem(w, req)

		if w.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
		}
	})

	t.Run("bad order id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodDelete, "/orders/abc/items/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc", "item_id": "1"})
		w := httptest.NewRecorder()

		h.DeleteItem(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("bad item id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodDelete, "/orders/1/items/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1", "item_id": "abc"})
		w := httptest.NewRecorder()

		h.DeleteItem(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("item not found", func(t *testing.T) {
		order := newTestOrder(t)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodDelete, "/orders/1/items/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1", "item_id": "1"})
		w := httptest.NewRecorder()

		h.DeleteItem(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestOrderHandler_PayOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)
		addTestItem(t, order)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				if updated.Status() != domain_order.StatusPaid {
					t.Fatalf("updated.Status() = %v, want %v", updated.Status(), domain_order.StatusPaid)
				}
				return nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/pay", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.PayOrder(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/abc/pay", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.PayOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("empty order returns conflict", func(t *testing.T) {
		order := newTestOrder(t)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/pay", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.PayOrder(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})
}

func TestOrderHandler_ShipOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := domain_order.NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 10),
			domain_order.StatusPaid,
			time.Now(),
		)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				if updated.Status() != domain_order.StatusShipped {
					t.Fatalf("updated.Status() = %v, want %v", updated.Status(), domain_order.StatusShipped)
				}
				return nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/ship", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ShipOrder(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/abc/ship", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.ShipOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("cannot ship returns conflict", func(t *testing.T) {
		order := newTestOrder(t)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/ship", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ShipOrder(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})
}

func TestOrderHandler_CancelOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		order := newTestOrder(t)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
			updateFn: func(ctx context.Context, updated *domain_order.Order) error {
				if updated.Status() != domain_order.StatusCancelled {
					t.Fatalf("updated.Status() = %v, want %v", updated.Status(), domain_order.StatusCancelled)
				}
				return nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/cancel", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.CancelOrder(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newOrderHandler(&productServiceMock{}, &orderRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/orders/abc/cancel", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.CancelOrder(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("cannot cancel returns conflict", func(t *testing.T) {
		order := domain_order.NewOrder(
			mustOrderID(t, 1),
			mustCustomerID(t, 10),
			domain_order.StatusShipped,
			time.Now(),
		)

		repo := &orderRepoMock{
			getFn: func(ctx context.Context, orderID domain_order.OrderID) (*domain_order.Order, error) {
				return order, nil
			},
		}

		h := newOrderHandler(&productServiceMock{}, repo)

		req := httptest.NewRequest(http.MethodPost, "/orders/1/cancel", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.CancelOrder(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})
}
