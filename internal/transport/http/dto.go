package http

import (
	"time"
)

type OrderItemDTO struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type OrderDTO struct {
	ID         int64          `json:"id"`
	CustomerID int64          `json:"customer_id"`
	Status     string         `json:"status"`
	CreatedAt  time.Time      `json:"created_at"`
	Items      []OrderItemDTO `json:"items"`
	Total      float64        `json:"total"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateOrderRequestDTO struct {
	CustomerID int64 `json:"customer_id"`
}

type CreateOrderResponseDTO struct {
	OrderID int64 `json:"order_id"`
}
