package product

import "errors"

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrInvalidProductID     = errors.New("invalid product id")
	ErrInvalidProductName   = errors.New("invalid product name")
	ErrInvalidPrice         = errors.New("invalid price")
	ErrInvalidStock         = errors.New("invalid stock")
	ErrInactiveProduct      = errors.New("product is inactive")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrProductAlreadyActive = errors.New("product already active")
)
