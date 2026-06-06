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

type ProductListParams struct {
	cursor    ProductID
	hasCursor bool

	limit int64
}

const defaultProductLimit int64 = 20
const maxProductLimit int64 = 100

func NewProductListParams(cursor *int64, limit *int64) (ProductListParams, error) {
	hasCursor := false
	var cursorValue ProductID = 0
	if cursor != nil {
		hasCursor = true
		productID, err := NewProductID(*cursor)
		if err != nil {
			return ProductListParams{}, ErrInvalidProductCursor
		}
		cursorValue = productID
	}
	if limit == nil {
		return ProductListParams{
			cursor:    cursorValue,
			hasCursor: hasCursor,
			limit:     defaultProductLimit,
		}, nil
	}
	if *limit <= 0 {
		return ProductListParams{}, ErrInvalidProductLimit
	}
	if *limit > maxProductLimit {
		return ProductListParams{
			cursor:    cursorValue,
			hasCursor: hasCursor,
			limit:     maxProductLimit,
		}, nil
	}
	return ProductListParams{
		cursor:    cursorValue,
		hasCursor: hasCursor,
		limit:     *limit,
	}, nil
}

func (p *ProductListParams) GetCursor() ProductID {
	return p.cursor
}

func (p *ProductListParams) GetLimit() int64 {
	return p.limit
}

func (p *ProductListParams) HasCursor() bool {
	return p.hasCursor
}
