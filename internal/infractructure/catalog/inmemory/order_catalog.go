// Package inmemory ...
package inmemory

import do "Order-Management-System/internal/domain/order"

type InMemoryOrderCatalog struct {
	data map[string]do.Product
}

func NewInMemoryOrderCatalog() *InMemoryOrderCatalog {
	return &InMemoryOrderCatalog{
		data: make(map[string]do.Product),
	}
}

func (c *InMemoryOrderCatalog) GetProduct(name string) (do.Product, error) {
	product, ok := c.data[name]
	if !ok {
		return do.Product{}, do.ErrUnknownOrderItem
	}
	return product, nil
}

func (c *InMemoryOrderCatalog) AddProduct(product do.Product) {
	c.data[product.Name] = product
}
