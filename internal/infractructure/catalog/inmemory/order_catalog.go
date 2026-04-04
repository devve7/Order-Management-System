// Package inmemory ...
package inmemory

import (
	domain_order "Order-Management-System/internal/domain/order"
	"sync"
)

type InMemoryOrderCatalog struct {
	data map[string]domain_order.Product
	mtx  sync.RWMutex
}

func NewInMemoryOrderCatalog() *InMemoryOrderCatalog {
	return &InMemoryOrderCatalog{
		data: make(map[string]domain_order.Product),
	}
}

func (c *InMemoryOrderCatalog) GetProduct(name string) (domain_order.Product, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	product, ok := c.data[name]
	if !ok {
		return domain_order.Product{}, domain_order.ErrProductNotFound
	}
	return product, nil
}

func (c *InMemoryOrderCatalog) AddProduct(product domain_order.Product) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[product.Name] = product
}
