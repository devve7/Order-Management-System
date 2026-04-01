// Package http ...
package http

import (
	application_order "Order-Management-System/internal/application/order"
	domain_order "Order-Management-System/internal/domain/order"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	var req CustomerIDDTO
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

	resp := OrderIDDTO{
		OrderID: int64(orderID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write create order response: %v", err)
	}
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	orderID, err := domain_order.NewOrderID(reqOrderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	order, err := h.usecase.GetOrder(orderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	resp := ToOrderDTO(&order)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write get order response: %v", err)
	}
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.usecase.GetOrders()
	if err != nil {
		h.writeError(w, err, mapError(err))
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
		log.Printf("failed to write get order response: %v", err)
	}
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	var req NameDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		h.writeError(w, errors.New("name is required"), http.StatusBadRequest)
		return
	}
	orderID, err := domain_order.NewOrderID(reqOrderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}

	orderItemID, err := h.usecase.AddItem(orderID, req.Name)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	resp := OrderItemIDDTO{
		OrderItemID: int64(orderItemID),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to write add item response: %v", err)
	}
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	reqOrderItemIDString := mux.Vars(r)["item_id"]
	reqOrderItemID, err := strconv.ParseInt(reqOrderItemIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	orderID, err := domain_order.NewOrderID(reqOrderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	orderItemID, err := domain_order.NewItemID(reqOrderItemID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}

	err = h.usecase.RemoveItem(orderID, orderItemID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PayOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	orderID, err := domain_order.NewOrderID(reqOrderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	err = h.usecase.PayOrder(orderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) ShipOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	orderID, err := domain_order.NewOrderID(reqOrderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	err = h.usecase.ShipOrder(orderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	orderID, err := domain_order.NewOrderID(reqOrderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	err = h.usecase.CancelOrder(orderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
