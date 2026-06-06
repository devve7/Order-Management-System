// Package order ...
package order

import (
	application_order "Order-Management-System/internal/application/order"
	domain_order "Order-Management-System/internal/domain/order"
	domain_product "Order-Management-System/internal/domain/product"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	usecase *application_order.UseCase
	logger  *logrus.Logger
}

func NewOrderHandler(usecase *application_order.UseCase, logger *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func mapError(err error) int {
	switch {
	case errors.Is(err, domain_order.ErrInvalidID),
		errors.Is(err, domain_order.ErrInvalidPrice),
		errors.Is(err, domain_order.ErrInvalidProductName),
		errors.Is(err, domain_order.ErrInvalidQuantity),
		errors.Is(err, domain_order.ErrInvalidStatus),
		errors.Is(err, domain_product.ErrInvalidPrice),
		errors.Is(err, domain_product.ErrInvalidProductID),
		errors.Is(err, domain_product.ErrInvalidProductName),
		errors.Is(err, domain_product.ErrInvalidStock),
		errors.Is(err, domain_order.ErrInvalidOrderCursor),
		errors.Is(err, domain_order.ErrInvalidOrderLimit):
		return http.StatusBadRequest

	case errors.Is(err, domain_order.ErrOrderNotFound),
		errors.Is(err, domain_order.ErrOrderItemNotFound),
		errors.Is(err, domain_product.ErrProductNotFound):
		return http.StatusNotFound

	case errors.Is(err, domain_order.ErrCannotPay),
		errors.Is(err, domain_order.ErrCannotShip),
		errors.Is(err, domain_order.ErrCannotCancel),
		errors.Is(err, domain_order.ErrOrderPaid),
		errors.Is(err, domain_order.ErrOrderCancelled),
		errors.Is(err, domain_order.ErrOrderShipped),
		errors.Is(err, domain_order.ErrOrderEmpty),
		errors.Is(err, domain_order.ErrConcurrentUpdate),
		errors.Is(err, domain_product.ErrInactiveProduct),
		errors.Is(err, domain_product.ErrInsufficientStock),
		errors.Is(err, domain_product.ErrProductAlreadyActive):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}

func (h *OrderHandler) writeError(w http.ResponseWriter, r *http.Request, err error, status int) {
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

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CustomerIDDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}

	if req.CustomerID == nil {
		h.writeError(w, r, errors.New("customer_id is required"), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	orderID, err := h.usecase.CreateOrder(ctx, *req.CustomerID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}

	resp := OrderIDDTO{
		OrderID: int64(orderID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to write create order response")
	}
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	order, err := h.usecase.GetOrder(ctx, reqOrderID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	resp := ToOrderDTO(&order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.logger.WithError(err).Error("failed to write get order response")
	}
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orders, err := h.usecase.GetOrders(ctx)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	ordersDTO := make([]OrderDTO, 0, len(orders))
	for _, order := range orders {
		orderDTO := ToOrderDTO(&order)
		ordersDTO = append(ordersDTO, orderDTO)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(ordersDTO); err != nil {
		h.logger.WithError(err).Error("failed to write get orders response")
	}
}

func (h *OrderHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req NewOrderItemDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	if req.ProductID == nil {
		h.writeError(w, r, errors.New("product_id is required"), http.StatusBadRequest)
		return
	}
	if req.Quantity == nil {
		h.writeError(w, r, errors.New("quantity is required"), http.StatusBadRequest)
		return
	}
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.usecase.AddItem(ctx, reqOrderID, *req.ProductID, *req.Quantity)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	reqOrderItemIDString := mux.Vars(r)["item_id"]
	reqOrderItemID, err := strconv.ParseInt(reqOrderItemIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.usecase.RemoveItem(ctx, reqOrderID, reqOrderItemID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *OrderHandler) PayOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.usecase.PayOrder(ctx, reqOrderID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) ShipOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.usecase.ShipOrder(ctx, reqOrderID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, r, err, http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err = h.usecase.CancelOrder(ctx, reqOrderID)
	if err != nil {
		h.writeError(w, r, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
