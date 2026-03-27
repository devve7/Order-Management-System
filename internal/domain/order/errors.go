package order

import "errors"

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrOrderEmpty        = errors.New("order is empry")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrCannotShip        = errors.New("cannot ship")
	ErrCannotPay         = errors.New("cannot pay")
	ErrCannotCancel      = errors.New("cannot cancel")
	ErrInvalidPrice      = errors.New("invalid price")
	ErrEmptyID           = errors.New("id is empty")
	ErrOrderShipped      = errors.New("order is shipped")
	ErrOrderCancelled    = errors.New("order is cancelled")
	ErrUnknownOrderItem  = errors.New("unknown order item")
	ErrOrderItemNotFound = errors.New("order item not found")
)
