package order

import (
	domain_order "Order-Management-System/internal/domain/order"
	"time"
)

type OrderItemDTO struct {
	ID        int64
	ProductID int64
	Name      string
	Price     int64
	Quantity  int64
}

type OrderDTO struct {
	ID         int64
	CustomerID int64
	Status     string
	CreatedAt  time.Time
	Items      []OrderItemDTO
	Total      int64
}

func ToOrderDTO(order *domain_order.Order) OrderDTO {
	items := make([]OrderItemDTO, 0, len(order.Items()))

	for _, item := range order.Items() {
		items = append(items, OrderItemDTO{
			ID:        int64(item.ID()),
			ProductID: int64(item.ProductID()),
			Name:      string(item.Name()),
			Price:     int64(item.Price()),
			Quantity:  int64(item.Quantity()),
		})
	}

	return OrderDTO{
		ID:         int64(order.ID()),
		CustomerID: int64(order.CustomerID()),
		Status:     string(order.Status()),
		CreatedAt:  order.CreatedAt(),
		Items:      items,
		Total:      int64(order.Total()),
	}
}

type OrderListParamsDTO struct {
	Cursor *int64
	Limit  *int64
}
