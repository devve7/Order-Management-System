package order

import "sync"

type OrderItemFactory struct {
	nextID  ItemID
	catalog ProductCatalog
	mtx     sync.Mutex
}

func NewOrderItemFactory(catalog ProductCatalog) *OrderItemFactory {
	return &OrderItemFactory{
		nextID:  1,
		catalog: catalog,
	}
}

func (f *OrderItemFactory) New(name string) (*OrderItem, error) {
	product, err := f.catalog.GetProduct(name)
	if err != nil {
		return nil, err
	}

	f.mtx.Lock()
	item := &OrderItem{
		id:    f.nextID,
		name:  product.Name,
		price: product.Price,
	}
	f.nextID++
	f.mtx.Unlock()
	return item, nil
}

type ProductCatalog interface {
	GetProduct(name string) (Product, error)
}

type Product struct {
	ID    ProductID
	Name  string
	Price Price
}
