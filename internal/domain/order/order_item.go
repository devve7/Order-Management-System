package order

type OrderItem struct {
	id        ItemID
	productID ProductID
	name      ItemName
	price     Price
	quantity  Quantity
}

func NewOrderItem(productID ProductID, name ItemName, price Price, quantity Quantity) (*OrderItem, error) {
	return &OrderItem{
		id:        0,
		productID: productID,
		name:      name,
		price:     price,
		quantity:  quantity,
	}, nil
}
func RestoreOrderItem(id ItemID, productID ProductID, name ItemName, price Price, quantity Quantity) (*OrderItem, error) {
	return &OrderItem{
		id:        id,
		productID: productID,
		name:      name,
		price:     price,
		quantity:  quantity,
	}, nil
}

func (i OrderItem) HasID() bool {
	return i.id > 0
}

func (i OrderItem) ID() ItemID {
	return i.id
}

func (i OrderItem) ProductID() ProductID {
	return i.productID
}

func (i OrderItem) Name() ItemName {
	return i.name
}

func (i OrderItem) Price() Price {
	return i.price
}

func (i OrderItem) Quantity() Quantity {
	return i.quantity
}
