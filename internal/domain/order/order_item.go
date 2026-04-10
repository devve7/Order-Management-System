package order

type OrderItem struct {
	id        ItemID
	productID ProductID
	name      string
	price     Price
	quantity  int32
}

func (i OrderItem) ID() ItemID {
	return i.id
}

func (i OrderItem) ProductID() ProductID {
	return i.productID
}

func (i OrderItem) Name() string {
	return i.name
}

func (i OrderItem) Price() Price {
	return i.price
}

func (i OrderItem) Quantity() int32 {
	return i.quantity
}
