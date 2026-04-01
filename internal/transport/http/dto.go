package http

import (
	application_order "Order-Management-System/internal/application/order"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

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

func ToOrderDTO(orderDTO *application_order.OrderDTO) OrderDTO {
	items := make([]OrderItemDTO, 0, len(orderDTO.Items))

	for _, item := range orderDTO.Items {
		items = append(items, OrderItemDTO{
			ID:    int64(item.ID),
			Name:  item.Name,
			Price: float64(item.Price),
		})
	}

	return OrderDTO{
		ID:         int64(orderDTO.ID),
		CustomerID: int64(orderDTO.CustomerID),
		Status:     string(orderDTO.Status),
		CreatedAt:  orderDTO.CreatedAt,
		Items:      items,
		Total:      float64(orderDTO.Total),
	}
}

type CustomerIDDTO struct {
	CustomerID *int64 `json:"customer_id"`
}

type OrderIDDTO struct {
	OrderID int64 `json:"order_id"`
}

type NameDTO struct {
	Name string `json:"name"`
}

type OrderItemIDDTO struct {
	OrderItemID int64 `json:"order_item_id"`
}
