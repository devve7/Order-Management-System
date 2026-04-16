// Package product ...
package product

import (
	application_product "Order-Management-System/internal/application/product"
	domain_product "Order-Management-System/internal/domain/product"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	usecase *application_product.UseCase
	logger  *logrus.Logger
}

func NewProductHandler(usecase *application_product.UseCase, logger *logrus.Logger) *ProductHandler {
	return &ProductHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func mapError(err error) int {
	switch {
	case errors.Is(err, domain_product.ErrInvalidPrice),
		errors.Is(err, domain_product.ErrInvalidProductID),
		errors.Is(err, domain_product.ErrInvalidProductName),
		errors.Is(err, domain_product.ErrInvalidStock):
		return http.StatusBadRequest

	case errors.Is(err, domain_product.ErrProductNotFound):
		return http.StatusNotFound

	case errors.Is(err, domain_product.ErrInactiveProduct),
		errors.Is(err, domain_product.ErrInsufficientStock),
		errors.Is(err, domain_product.ErrProductAlreadyActive):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}

func (h *ProductHandler) writeError(w http.ResponseWriter, r *http.Request, err error, status int) {
	entry := h.logger.WithFields(logrus.Fields{
		"method":      r.Method,
		"uri":         r.RequestURI,
		"remote_addr": r.RemoteAddr,
		"user_agent":  r.UserAgent(),
		"status":      status,
	}).WithError(err)

	if status >= 500 {
		entry.Error("request failed")
	} else {
		entry.Warn("request failed")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: err.Error(),
	}

	if status == 500 {
		resp.Error = "internal server error"
	}

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		h.logger.WithError(err).Error("failed to write error response")
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		h.writeError(w, r, errors.New("name is required"), http.StatusBadRequest)
		return
	}
	if req.Price == nil {
		h.writeError(w, r, errors.New("price_cents is required"), http.StatusBadRequest)
		return
	}
	if req.Stock == nil {
		h.writeError(w, r, errors.New("stock is required"), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	id, err := h.usecase.CreateProduct(ctx, req.Name, *req.Price, *req.Stock)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
	}

	resp := ProductIDDTO{
		ProductID: id,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to write create product response")
	}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	reqProductIDString := mux.Vars(r)["id"]
	reqProductID, err := strconv.ParseInt(reqProductIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	product, err := h.usecase.GetProduct(ctx, reqProductID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	resp := ToProductDTO(product)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to write get product response")
	}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.usecase.GetProducts(ctx)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	productsDTO := make([]ProductDTO, 0, len(products))
	for _, product := range products {
		productDTO := ToProductDTO(product)
		productsDTO = append(productsDTO, productDTO)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(productsDTO); err != nil {
		h.logger.WithError(err).Error("failed to write get orders response")
	}
}

func (h *ProductHandler) ChangePrice(w http.ResponseWriter, r *http.Request) {
	reqProductIDString := mux.Vars(r)["id"]
	reqProductID, err := strconv.ParseInt(reqProductIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	var req PriceDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.usecase.ChangePrice(ctx, reqProductID, req.Price)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) AddStock(w http.ResponseWriter, r *http.Request) {
	reqProductIDString := mux.Vars(r)["id"]
	reqProductID, err := strconv.ParseInt(reqProductIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	var req StockAmountDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.usecase.AddStock(ctx, reqProductID, req.StockAmount)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) RemoveStock(w http.ResponseWriter, r *http.Request) {
	reqProductIDString := mux.Vars(r)["id"]
	reqProductID, err := strconv.ParseInt(reqProductIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	var req StockAmountDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.usecase.RemoveStock(ctx, reqProductID, req.StockAmount)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) ActivateProduct(w http.ResponseWriter, r *http.Request) {
	reqProductIDString := mux.Vars(r)["id"]
	reqProductID, err := strconv.ParseInt(reqProductIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.usecase.ActivateProduct(ctx, reqProductID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *ProductHandler) DeactivateProduct(w http.ResponseWriter, r *http.Request) {
	reqProductIDString := mux.Vars(r)["id"]
	reqProductID, err := strconv.ParseInt(reqProductIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.usecase.DeactivateProduct(ctx, reqProductID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
