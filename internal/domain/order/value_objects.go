package order

type OrderID int64

func NewOrderID(id int64) (OrderID, error) {
	if id < 0 {
		return 0, ErrInvalidID
	}
	return OrderID(id), nil
}

type CustomerID int64

func NewCustomerID(id int64) (CustomerID, error) {
	if id < 0 {
		return 0, ErrInvalidID
	}
	return CustomerID(id), nil
}

type ItemID int64

func NewItemID(id int64) (ItemID, error) {
	if id < 0 {
		return 0, ErrInvalidID
	}
	return ItemID(id), nil
}

type ProductID int64

func NewProductID(id int64) (ProductID, error) {
	if id < 0 {
		return 0, ErrInvalidID
	}
	return ProductID(id), nil
}

type Price int64

func NewPrice(amount int64) (Price, error) {
	if amount < 0 {
		return 0, ErrInvalidPrice
	}
	return Price(amount), nil
}

type Quantity int64

func NewQuantity(q int64) (Quantity, error) {
	if q <= 0 {
		return 0, ErrInvalidQuantity
	}
	return Quantity(q), nil
}

type ProductName string

func NewProductName(s string) (ProductName, error) {
	if s == "" {
		return "", ErrInvalidProductName
	}
	return ProductName(s), nil
}

type OrderVersion int64

func NewOrderVersion(v int64) (OrderVersion, error) {
	if v <= 0 {
		return 0, ErrInvalidOrderVersion
	}
	return OrderVersion(v), nil
}

type OrderListParams struct {
	customerID CustomerID

	cursor    OrderID
	hasCursor bool

	limit int64
}

const defaultOrderLimit int64 = 20
const maxOrderLimit int64 = 100

func NewOrderListParams(customerID CustomerID, cursor *int64, limit *int64) (OrderListParams, error) {
	hasCursor := false
	var cursorValue OrderID = 0
	if cursor != nil {
		hasCursor = true
		orderID, err := NewOrderID(*cursor)
		if err != nil {
			return OrderListParams{}, ErrInvalidOrderCursor
		}
		cursorValue = orderID
	}
	if limit == nil {
		return OrderListParams{
			customerID: customerID,
			cursor:     cursorValue,
			hasCursor:  hasCursor,
			limit:      defaultOrderLimit,
		}, nil
	}
	if *limit <= 0 {
		return OrderListParams{}, ErrInvalidOrderLimit
	}
	if *limit > maxOrderLimit {
		return OrderListParams{
			customerID: customerID,
			cursor:     cursorValue,
			hasCursor:  hasCursor,
			limit:      maxOrderLimit,
		}, nil
	}
	return OrderListParams{
		customerID: customerID,
		cursor:     cursorValue,
		hasCursor:  hasCursor,
		limit:      *limit,
	}, nil
}

func (p *OrderListParams) Cursor() OrderID {
	return p.cursor
}

func (p *OrderListParams) Limit() int64 {
	return p.limit
}

func (p *OrderListParams) HasCursor() bool {
	return p.hasCursor
}

func (p *OrderListParams) CustomerID() CustomerID {
	return p.customerID
}
