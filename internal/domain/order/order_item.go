package order

import "sync"

type OrderItem struct {
	id    ItemID
	name  string
	price Price
}

func (i *OrderItem) GetID() ItemID {
	return i.id
}

type OrderItemFactory struct {
	nextID ItemID
	prices map[string]Price
	mtx    sync.Mutex
}

func NewOrderItemFactory() *OrderItemFactory {
	return &OrderItemFactory{
		nextID: 1,
		prices: make(map[string]Price),
	}
}

func (f *OrderItemFactory) AddItem(name string, price Price) {
	f.mtx.Lock()
	f.prices[name] = price
	f.mtx.Unlock()
}

func (f *OrderItemFactory) NewOrderItem(name string) (*OrderItem, error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	price, ok := f.prices[name]
	if !ok {
		return &OrderItem{}, ErrUnknownOrderItem
	}

	id := f.nextID
	f.nextID = NewItemID(uint64(id) + 1)

	return &OrderItem{
		id:    id,
		name:  name,
		price: price,
	}, nil
}
