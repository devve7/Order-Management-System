package order

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")
	ErrOrderEmpty    = errors.New("order is empry")
	ErrInvalidStatus = errors.New("invalid status")
	ErrCannotShip    = errors.New("cannot ship")
	ErrCannotCancel  = errors.New("cannot cancel")
	ErrInvalidPrice  = errors.New("invalid price")
	ErrEmptyID       = errors.New("id is empty")
)
