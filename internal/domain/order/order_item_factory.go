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

func (f *OrderItemFactory) New(productID ProductID, quantity int64) (*OrderItem, error) {
	product, err := f.catalog.GetProduct(productID)
	if err != nil {
		return nil, err
	}

	f.mtx.Lock()
	item := &OrderItem{
		id:        f.nextID,
		productID: product.ID,
		name:      product.Name,
		price:     product.Price,
		quantity:  quantity,
	}
	f.nextID++
	f.mtx.Unlock()
	return item, nil
}

type ProductCatalog interface {
	GetProduct(productID ProductID) (Product, error)
}

type Product struct {
	ID    ProductID
	Name  string
	Price Price
}
