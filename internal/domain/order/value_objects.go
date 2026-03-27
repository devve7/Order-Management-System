package order

type OrderID string

func NewOrderID(id string) (OrderID, error) {
	if id == "" {
		return "", ErrEmptyID
	}
	return OrderID(id), nil
}

type CustomerID string

func NewCustomerID(id string) (CustomerID, error) {
	if id == "" {
		return "", ErrEmptyID
	}
	return CustomerID(id), nil
}

type ItemID string

func NewItemID(id string) (ItemID, error) {
	if id == "" {
		return "", ErrEmptyID
	}
	return ItemID(id), nil
}

type Price float64

func NewPrice(amount float64) (Price, error) {
	if amount <= 0 {
		return 0, ErrInvalidPrice
	}
	return Price(amount), nil
}
