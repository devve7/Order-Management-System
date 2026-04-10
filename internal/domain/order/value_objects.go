package order

type OrderID int64

func NewOrderID(id int64) (OrderID, error) {
	if id < 0 {
		return OrderID(0), ErrInvalidID
	}
	return OrderID(id), nil
}

type CustomerID int64

func NewCustomerID(id int64) (CustomerID, error) {
	if id < 0 {
		return CustomerID(0), ErrInvalidID
	}
	return CustomerID(id), nil
}

type ItemID int64

func NewItemID(id int64) (ItemID, error) {
	if id < 0 {
		return ItemID(0), ErrInvalidID
	}
	return ItemID(id), nil
}

type ProductID int64

func NewProductID(id int64) (ProductID, error) {
	if id < 0 {
		return ProductID(0), ErrInvalidID
	}
	return ProductID(id), nil
}

type Price float64

func NewPrice(amount float64) (Price, error) {
	if amount < 0 {
		return 0, ErrInvalidPrice
	}
	return Price(amount), nil
}

type Quantity int64

func NewQuantity(q int64) (Quantity, error) {
	if q < 0 {
		return 0, ErrInvalidQuantity
	}
	return Quantity(q), nil
}

type ItemName string

func NewItemName(s string) (ItemName, error) {
	if s == "" {
		return "", ErrInvalidItemName
	}
	return ItemName(s), nil
}
