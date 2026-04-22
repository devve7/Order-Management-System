package product

import (
	application_product "Order-Management-System/internal/application/product"
	domain_product "Order-Management-System/internal/domain/product"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type productRepoMock struct {
	createFn func(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error)
	getFn    func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error)
	getAllFn func(ctx context.Context) ([]*domain_product.Product, error)
	updateFn func(ctx context.Context, product *domain_product.Product) error
}

func (m *productRepoMock) Create(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
	if m.createFn != nil {
		return m.createFn(ctx, product)
	}
	return 0, nil
}

func (m *productRepoMock) Get(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
	if m.getFn != nil {
		return m.getFn(ctx, id)
	}
	return nil, nil
}

func (m *productRepoMock) GetAll(ctx context.Context) ([]*domain_product.Product, error) {
	if m.getAllFn != nil {
		return m.getAllFn(ctx)
	}
	return nil, nil
}

func (m *productRepoMock) Update(ctx context.Context, product *domain_product.Product) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, product)
	}
	return nil
}

func newProductHandler(repo domain_product.Repository) *ProductHandler {
	usecase := application_product.NewUseCase(repo)
	logger := logrus.New()
	logger.SetOutput(bytes.NewBuffer(nil))
	return NewProductHandler(usecase, logger)
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

func decodeErrorResponse(t *testing.T, rr *httptest.ResponseRecorder) ErrorResponse {
	t.Helper()
	var resp ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	return resp
}

func TestProductHandler_CreateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &productRepoMock{
			createFn: func(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
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

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(
			`{"name":"iPhone","price_cents":100000,"stock":10}`,
		))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
		}
		if got := w.Header().Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want %q", got, "application/json")
		}

		var resp ProductIDDTO
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if resp.ProductID != 42 {
			t.Fatalf("product_id = %d, want 42", resp.ProductID)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(`{"name":`))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(
			`{"name":"iPhone","price_cents":100000,"stock":10,"extra":"x"}`,
		))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("missing name", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(
			`{"price_cents":100000,"stock":10}`,
		))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		resp := decodeErrorResponse(t, w)
		if resp.Error != "name is required" {
			t.Fatalf("error = %q, want %q", resp.Error, "name is required")
		}
	})

	t.Run("missing price", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(
			`{"name":"iPhone","stock":10}`,
		))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		resp := decodeErrorResponse(t, w)
		if resp.Error != "price_cents is required" {
			t.Fatalf("error = %q, want %q", resp.Error, "price_cents is required")
		}
	})

	t.Run("missing stock", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(
			`{"name":"iPhone","price_cents":100000}`,
		))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		resp := decodeErrorResponse(t, w)
		if resp.Error != "stock is required" {
			t.Fatalf("error = %q, want %q", resp.Error, "stock is required")
		}
	})

	t.Run("invalid price from usecase", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(
			`{"name":"iPhone","price_cents":-1,"stock":10}`,
		))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("repo create error becomes 500", func(t *testing.T) {
		repo := &productRepoMock{
			createFn: func(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
				return 0, errors.New("db exploded")
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBufferString(
			`{"name":"iPhone","price_cents":100000,"stock":10}`,
		))
		w := httptest.NewRecorder()

		h.CreateProduct(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		resp := decodeErrorResponse(t, w)
		if resp.Error != "internal server error" {
			t.Fatalf("error = %q, want %q", resp.Error, "internal server error")
		}
	})
}

func TestProductHandler_GetProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 7),
			mustProductName(t, "MacBook"),
			mustPrice(t, 200000),
			mustStock(t, 5),
			true,
		)

		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				if id != 7 {
					t.Fatalf("id = %v, want %v", id, domain_product.ProductID(7))
				}
				return product, nil
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/products/7", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "7"})
		w := httptest.NewRecorder()

		h.GetProduct(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if got := w.Header().Get("Content-Type"); got != "application/json" {
			t.Fatalf("Content-Type = %q, want %q", got, "application/json")
		}

		var resp ProductDTO
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if resp.ID != 7 || resp.Name != "MacBook" || resp.Price != 200000 || resp.Stock != 5 || resp.Active != true {
			t.Fatalf("unexpected dto: %+v", resp)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("zero id maps to bad request", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodGet, "/products/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		h.GetProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, domain_product.ErrProductNotFound
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/products/7", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "7"})
		w := httptest.NewRecorder()

		h.GetProduct(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("unexpected repo error becomes 500", func(t *testing.T) {
		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, errors.New("db failed")
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/products/7", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "7"})
		w := httptest.NewRecorder()

		h.GetProduct(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestProductHandler_GetProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		products := []*domain_product.Product{
			domain_product.RestoreProduct(
				mustProductID(t, 1),
				mustProductName(t, "iPhone"),
				mustPrice(t, 100000),
				mustStock(t, 10),
				true,
			),
			domain_product.RestoreProduct(
				mustProductID(t, 2),
				mustProductName(t, "MacBook"),
				mustPrice(t, 200000),
				mustStock(t, 5),
				false,
			),
		}

		repo := &productRepoMock{
			getAllFn: func(ctx context.Context) ([]*domain_product.Product, error) {
				return products, nil
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()

		h.GetProducts(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp []ProductDTO
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatalf("decode response: %v", err)
		}

		if len(resp) != 2 {
			t.Fatalf("len(resp) = %d, want 2", len(resp))
		}
		if resp[0].ID != 1 || resp[1].ID != 2 {
			t.Fatalf("ids = [%d %d], want [1 2]", resp[0].ID, resp[1].ID)
		}
	})

	t.Run("repo error becomes 500", func(t *testing.T) {
		repo := &productRepoMock{
			getAllFn: func(ctx context.Context) ([]*domain_product.Product, error) {
				return nil, errors.New("db failed")
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()

		h.GetProducts(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestProductHandler_ChangePrice(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				if updated.Price() != 120000 {
					t.Fatalf("updated.Price() = %v, want %v", updated.Price(), domain_product.Price(120000))
				}
				return nil
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPatch, "/products/1/price", bytes.NewBufferString(`{"price":120000}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ChangePrice(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPatch, "/products/abc/price", bytes.NewBufferString(`{"price":120000}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.ChangePrice(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPatch, "/products/1/price", bytes.NewBufferString(`{"price":`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ChangePrice(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPatch, "/products/1/price", bytes.NewBufferString(`{"price":120000,"extra":1}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ChangePrice(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid price from usecase", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPatch, "/products/1/price", bytes.NewBufferString(`{"price":-1}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ChangePrice(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, domain_product.ErrProductNotFound
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPatch, "/products/1/price", bytes.NewBufferString(`{"price":120000}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ChangePrice(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("update error becomes 500", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				return errors.New("db failed")
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPatch, "/products/1/price", bytes.NewBufferString(`{"price":120000}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ChangePrice(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestProductHandler_AddStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &productRepoMock{
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

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/add", bytes.NewBufferString(`{"amount":5}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddStock(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/abc/stock/add", bytes.NewBufferString(`{"amount":5}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.AddStock(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/add", bytes.NewBufferString(`{"amount":`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddStock(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/add", bytes.NewBufferString(`{"amount":5,"extra":1}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddStock(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid amount from usecase", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/add", bytes.NewBufferString(`{"amount":-1}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.AddStock(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

func TestProductHandler_RemoveStock(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
			updateFn: func(ctx context.Context, updated *domain_product.Product) error {
				if updated.Stock() != 6 {
					t.Fatalf("updated.Stock() = %v, want %v", updated.Stock(), domain_product.Stock(6))
				}
				return nil
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/remove", bytes.NewBufferString(`{"amount":4}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.RemoveStock(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/abc/stock/remove", bytes.NewBufferString(`{"amount":4}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.RemoveStock(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/remove", bytes.NewBufferString(`{"amount":`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.RemoveStock(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/remove", bytes.NewBufferString(`{"amount":4,"extra":1}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.RemoveStock(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("insufficient stock returns conflict", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/remove", bytes.NewBufferString(`{"amount":11}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.RemoveStock(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})

	t.Run("not found", func(t *testing.T) {
		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return nil, domain_product.ErrProductNotFound
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock/remove", bytes.NewBufferString(`{"amount":4}`))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.RemoveStock(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestProductHandler_ActivateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		repo := &productRepoMock{
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

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/activate", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ActivateProduct(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/abc/activate", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.ActivateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("already active returns conflict", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/activate", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.ActivateProduct(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})
}

func TestProductHandler_DeactivateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			true,
		)

		repo := &productRepoMock{
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

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/deactivate", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.DeactivateProduct(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("bad id", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})

		req := httptest.NewRequest(http.MethodPost, "/products/abc/deactivate", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.DeactivateProduct(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("already inactive returns conflict", func(t *testing.T) {
		product := domain_product.RestoreProduct(
			mustProductID(t, 1),
			mustProductName(t, "iPhone"),
			mustPrice(t, 100000),
			mustStock(t, 10),
			false,
		)

		repo := &productRepoMock{
			getFn: func(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
				return product, nil
			},
		}

		h := newProductHandler(repo)

		req := httptest.NewRequest(http.MethodPost, "/products/1/deactivate", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.DeactivateProduct(w, req)

		if w.Code != http.StatusConflict {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})
}

func TestProductHandler_mapError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"invalid product id", domain_product.ErrInvalidProductID, http.StatusBadRequest},
		{"invalid product name", domain_product.ErrInvalidProductName, http.StatusBadRequest},
		{"invalid price", domain_product.ErrInvalidPrice, http.StatusBadRequest},
		{"invalid stock", domain_product.ErrInvalidStock, http.StatusBadRequest},
		{"not found", domain_product.ErrProductNotFound, http.StatusNotFound},
		{"inactive", domain_product.ErrInactiveProduct, http.StatusConflict},
		{"insufficient stock", domain_product.ErrInsufficientStock, http.StatusConflict},
		{"already active", domain_product.ErrProductAlreadyActive, http.StatusConflict},
		{"unknown", errors.New("boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapError(tt.err)
			if got != tt.want {
				t.Fatalf("mapError(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}

func TestProductHandler_writeError(t *testing.T) {
	t.Run("client error keeps original message", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})
		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		w := httptest.NewRecorder()

		h.writeError(w, req, domain_product.ErrInvalidProductID, http.StatusBadRequest)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		resp := decodeErrorResponse(t, w)
		if resp.Error != domain_product.ErrInvalidProductID.Error() {
			t.Fatalf("error = %q, want %q", resp.Error, domain_product.ErrInvalidProductID.Error())
		}
	})

	t.Run("server error hides internal message", func(t *testing.T) {
		h := newProductHandler(&productRepoMock{})
		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		w := httptest.NewRecorder()

		h.writeError(w, req, errors.New("sql is down"), http.StatusInternalServerError)

		if w.Code != http.StatusInternalServerError {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		resp := decodeErrorResponse(t, w)
		if resp.Error != "internal server error" {
			t.Fatalf("error = %q, want %q", resp.Error, "internal server error")
		}
	})
}
