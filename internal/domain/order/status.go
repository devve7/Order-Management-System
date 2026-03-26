package order

type OrderStatus string

const (
	StatusCreated   OrderStatus = "created"
	StatusPaid      OrderStatus = "paid"
	StatusShipped   OrderStatus = "shipped"
	StatusCancelled OrderStatus = "cancelled"
)

func (s OrderStatus) IsValid() bool {
	switch s {
	case StatusCancelled, StatusCreated, StatusPaid, StatusShipped:
		return true
	default:
		return false
	}
}

func NewOrderStatus(status string) (OrderStatus, error) {
	s := OrderStatus(status)
	if !s.IsValid() {
		return "", ErrInvalidStatus
	}
	return s, nil
}
