// Package http ...
package http

import (
	application_order "Order-Management-System/internal/application/order"
	domain_order "Order-Management-System/internal/domain/order"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	usecase *application_order.UseCase
}

func NewHandler(usecase *application_order.UseCase) *Handler {
	return &Handler{
		usecase: usecase,
	}
}

func mapError(err error) int {
	switch {
	case errors.Is(err, domain_order.ErrInvalidID),
		errors.Is(err, domain_order.ErrInvalidPrice),
		errors.Is(err, domain_order.ErrInvalidStatus):
		return http.StatusBadRequest

	case errors.Is(err, domain_order.ErrOrderNotFound),
		errors.Is(err, domain_order.ErrProductNotFound),
		errors.Is(err, domain_order.ErrOrderItemNotFound):
		return http.StatusNotFound

	case errors.Is(err, domain_order.ErrCannotPay),
		errors.Is(err, domain_order.ErrCannotShip),
		errors.Is(err, domain_order.ErrCannotCancel),
		errors.Is(err, domain_order.ErrOrderCancelled),
		errors.Is(err, domain_order.ErrOrderShipped),
		errors.Is(err, domain_order.ErrOrderEmpty):
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}

func (h *Handler) writeError(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := ErrorResponse{
		Error: err.Error(),
	}

	if status == 500 {
		resp.Error = "internal server error"
	}

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		log.Printf("failed to write error response: %v", encErr)
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequestDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}

	if req.CustomerID == nil {
		h.writeError(w, errors.New("customer_id is required"), http.StatusBadRequest)
		return
	}

	customerID, err := domain_order.NewCustomerID(*req.CustomerID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}

	orderID, err := h.usecase.CreateOrder(customerID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}

	resp := CreateOrderResponseDTO{
		OrderID: int64(orderID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write create order response: %v", err)
	}
}
