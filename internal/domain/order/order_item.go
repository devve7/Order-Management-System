package order

type OrderItem struct {
	id    ItemID
	name  string
	price Price
}

func (i OrderItem) ID() ItemID {
	return i.id
}

func (i OrderItem) Name() string {
	return i.name
}

func (i OrderItem) Price() Price {
	return i.price
}
