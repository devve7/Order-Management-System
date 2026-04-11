// Package http ...
package http

import (
	application_order "Order-Management-System/internal/application/order"
	application_product "Order-Management-System/internal/application/product"
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
		errors.Is(err, domain_order.ErrInvalidProductName),
		errors.Is(err, domain_order.ErrInvalidQuantity),
		errors.Is(err, domain_order.ErrInvalidStatus):
		return http.StatusBadRequest

	case errors.Is(err, domain_order.ErrOrderNotFound),
		errors.Is(err, domain_order.ErrOrderItemNotFound),
		errors.Is(err, application_product.ErrProductNotFound):
		return http.StatusNotFound

	case errors.Is(err, domain_order.ErrCannotPay),
		errors.Is(err, domain_order.ErrCannotShip),
		errors.Is(err, domain_order.ErrCannotCancel),
		errors.Is(err, domain_order.ErrOrderCancelled),
		errors.Is(err, domain_order.ErrOrderShipped),
		errors.Is(err, domain_order.ErrOrderEmpty),
		errors.Is(err, application_product.ErrInsufficientStock),
		errors.Is(err, domain_order.ErrConcurrentUpdate):
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
	ctx := r.Context()
	orderID, err := h.usecase.CreateOrder(ctx, customerID)
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
	ctx := r.Context()
	order, err := h.usecase.GetOrder(ctx, orderID)
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
	ctx := r.Context()
	orders, err := h.usecase.GetOrders(ctx)
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
	var req NewOrderItemDTO
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	if req.ProductID == nil {
		h.writeError(w, errors.New("product_id is required"), http.StatusBadRequest)
		return
	}
	if req.Quantity == nil {
		h.writeError(w, errors.New("quantity is required"), http.StatusBadRequest)
		return
	}
	reqOrderIDString := mux.Vars(r)["id"]
	reqOrderID, err := strconv.ParseInt(reqOrderIDString, 10, 64)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	productID, err := domain_order.NewProductID(*req.ProductID)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	quantity, err := domain_order.NewQuantity(*req.Quantity)
	if err != nil {
		h.writeError(w, err, http.StatusBadRequest)
		return
	}
	orderID, err := domain_order.NewOrderID(reqOrderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	ctx := r.Context()
	err = h.usecase.AddItem(ctx, orderID, productID, quantity)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
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
	ctx := r.Context()
	err = h.usecase.RemoveItem(ctx, orderID, orderItemID)
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
	ctx := r.Context()
	err = h.usecase.PayOrder(ctx, orderID)
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
	ctx := r.Context()
	err = h.usecase.ShipOrder(ctx, orderID)
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
	ctx := r.Context()
	err = h.usecase.CancelOrder(ctx, orderID)
	if err != nil {
		h.writeError(w, err, mapError(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
