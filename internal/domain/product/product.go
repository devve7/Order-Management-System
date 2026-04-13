// Package product ...
package product

type Product struct {
	id     ProductID
	name   ProductName
	price  Price
	stock  Stock
	active bool
}

func NewProduct(name ProductName, price Price, stock Stock) *Product {
	return &Product{
		id:     0,
		name:   name,
		price:  price,
		stock:  stock,
		active: true,
	}
}

func RestoreProduct(id ProductID, name ProductName, price Price, stock Stock, active bool) *Product {
	return &Product{
		id:     id,
		name:   name,
		price:  price,
		stock:  stock,
		active: active,
	}
}

func (p *Product) ChangePrice(price Price) {
	p.price = price
}

func (p *Product) AddStock(stock Stock) {
	p.stock += stock
}

func (p *Product) RemoveStock(stock Stock) error {
	if p.stock < stock {
		return ErrInsufficientStock
	}
	p.stock -= stock
	return nil
}

func (p *Product) Deactivate() error {
	if !p.active {
		return ErrInactiveProduct
	}
	p.active = false
	return nil
}

func (p *Product) Activate() error {
	if p.active {
		return ErrProductAlreadyActive
	}
	p.active = true
	return nil
}

func (p *Product) ID() ProductID {
	return p.id
}

func (p *Product) Name() ProductName {
	return p.name
}

func (p *Product) Price() Price {
	return p.price
}

func (p *Product) Stock() Stock {
	return p.stock
}

func (p *Product) Active() bool {
	return p.active
}

func (p *Product) HasID() bool {
	return p.id > 0
}
