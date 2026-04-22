package product

type ProductID int64

func NewProductID(id int64) (ProductID, error) {
	if id <= 0 {
		return 0, ErrInvalidProductID
	}
	return ProductID(id), nil
}

type ProductName string

func NewProductName(name string) (ProductName, error) {
	if name == "" {
		return "", ErrInvalidProductName
	}
	return ProductName(name), nil
}

type Price int64

func NewPrice(price int64) (Price, error) {
	if price < 0 {
		return 0, ErrInvalidPrice
	}
	return Price(price), nil
}

type Stock int64

func NewStock(stock int64) (Stock, error) {
	if stock < 0 {
		return 0, ErrInvalidStock
	}
	return Stock(stock), nil
}
