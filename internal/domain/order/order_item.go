package order

type OrderItem struct {
	id    ItemID
	name  string
	price Price
}

func (i OrderItem) ID() ItemID {
	return i.id
}
