// Package inmemory ...
package inmemory

import (
	do "Order-Management-System/internal/domain/order"
	"sync"
)

type InMemoryOrderCatalog struct {
	data map[string]do.Product
	mtx  sync.RWMutex
}

func NewInMemoryOrderCatalog() *InMemoryOrderCatalog {
	return &InMemoryOrderCatalog{
		data: make(map[string]do.Product),
	}
}

func (c *InMemoryOrderCatalog) GetProduct(name string) (do.Product, error) {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	product, ok := c.data[name]
	if !ok {
		return do.Product{}, do.ErrProductNotFound
	}
	return product, nil
}

func (c *InMemoryOrderCatalog) AddProduct(product do.Product) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.data[product.Name] = product
}
