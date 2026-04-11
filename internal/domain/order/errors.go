package order

import "errors"

var (
	ErrOrderNotFound  = errors.New("order not found")
	ErrOrderEmpty     = errors.New("order is empty")
	ErrOrderShipped   = errors.New("order is shipped")
	ErrOrderCancelled = errors.New("order is cancelled")
	ErrOrderPaid      = errors.New("order is paid")

	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidPrice        = errors.New("invalid price")
	ErrInvalidID           = errors.New("id cannot be negative")
	ErrInvalidQuantity     = errors.New("invalid quantity")
	ErrInvalidProductName  = errors.New("invalid product name")
	ErrInvalidOrderVersion = errors.New("invalid order version")

	ErrCannotShip   = errors.New("cannot ship")
	ErrCannotPay    = errors.New("cannot pay")
	ErrCannotCancel = errors.New("cannot cancel")

	ErrOrderItemNotFound = errors.New("order item not found")

	ErrConcurrentUpdate = errors.New("order version conflict")
)
