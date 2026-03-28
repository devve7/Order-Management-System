package order

type OrderItem struct {
	id    ItemID
	name  string
	price Price
}

func (i OrderItem) ID() ItemID {
	return i.id
}

type OrderItemFactory struct {
	catalog ProductCatalog
}

func NewOrderItemFactory(catalog ProductCatalog) *OrderItemFactory {
	return &OrderItemFactory{
		catalog: catalog,
	}
}

func (f *OrderItemFactory) New(name string) (*OrderItem, error) {
	product, err := f.catalog.GetProduct(name)
	if err != nil {
		return nil, err
	}

	return &OrderItem{
		id:    product.ID,
		name:  product.Name,
		price: product.Price,
	}, nil
}

type ProductCatalog interface {
	GetProduct(name string) (Product, error)
}

type Product struct {
	ID    ItemID
	Name  string
	Price Price
}
