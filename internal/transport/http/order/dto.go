package order

import (
	application_order "Order-Management-System/internal/application/order"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type OrderItemDTO struct {
	ID        int64  `json:"order_item_id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Quantity  int64  `json:"quantity"`
}

type OrderDTO struct {
	ID         int64          `json:"order_id"`
	CustomerID int64          `json:"customer_id"`
	Status     string         `json:"status"`
	CreatedAt  string         `json:"created_at"`
	Items      []OrderItemDTO `json:"items"`
	Total      int64          `json:"total"`
}

func ToOrderDTO(orderDTO *application_order.OrderDTO) OrderDTO {
	items := make([]OrderItemDTO, 0, len(orderDTO.Items))

	for _, item := range orderDTO.Items {
		items = append(items, OrderItemDTO{
			ID:        item.ID,
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Quantity:  item.Quantity,
		})
	}

	return OrderDTO{
		ID:         orderDTO.ID,
		CustomerID: orderDTO.CustomerID,
		Status:     orderDTO.Status,
		CreatedAt:  orderDTO.CreatedAt.Format("2006-01-02 15:04:05"),
		Items:      items,
		Total:      orderDTO.Total,
	}
}

type CustomerIDDTO struct {
	CustomerID *int64 `json:"customer_id"`
}

type OrderIDDTO struct {
	OrderID int64 `json:"order_id"`
}

type NewOrderItemDTO struct {
	ProductID *int64 `json:"product_id"`
	Quantity  *int64 `json:"quantity"`
}

type OrderListQuery struct {
	Cursor *int64
	Limit  *int64
}
