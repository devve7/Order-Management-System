package order

import (
	do "Order-Management-System/internal/domain/order"
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

func ToOrderDTO(order *do.Order) OrderDTO {
	items := make([]OrderItemDTO, 0, len(order.Items()))

	for _, item := range order.Items() {
		items = append(items, OrderItemDTO{
			ID:    int64(item.ID()),
			Name:  item.Name(),
			Price: float64(item.Price()),
		})
	}

	return OrderDTO{
		ID:         int64(order.ID()),
		CustomerID: int64(order.CustomerID()),
		Status:     string(order.Status()),
		CreatedAt:  order.CreatedAt(),
		Items:      items,
		Total:      float64(order.Total()),
	}
}
