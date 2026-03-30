// Package order ...
package order

import (
	"time"
)

type Order struct {
	id         OrderID
	customerID CustomerID
	items      []*OrderItem
	status     OrderStatus
	createdAt  time.Time
}

func NewOrder(id OrderID, customerID CustomerID) *Order {
	status := StatusCreated
	return &Order{
		id:         id,
		customerID: customerID,
		items:      make([]*OrderItem, 0),
		status:     status,
		createdAt:  time.Now(),
	}
}

func (o *Order) AddItem(item *OrderItem) error {
	if item == nil {
		return ErrOrderItemNotFound
	}
	if o.status == StatusCancelled {
		return ErrOrderCancelled
	}
	if o.status == StatusShipped {
		return ErrOrderShipped
	}
	o.items = append(o.items, item)

	return nil
}

func (o *Order) RemoveItem(id ItemID) error {
	if o.status == StatusCancelled {
		return ErrOrderCancelled
	}
	if o.status == StatusShipped {
		return ErrOrderShipped
	}
	if len(o.items) == 0 {
		return ErrOrderItemNotFound
	}
	for index, item := range o.items {
		if item.ID() == id {
			o.items = append(o.items[:index], o.items[index+1:]...)
			return nil
		}
	}
	return ErrOrderItemNotFound
}

func (o *Order) Pay() error {
	if len(o.items) == 0 {
		return ErrOrderEmpty
	}
	if o.status == StatusCreated {
		o.status = StatusPaid
		return nil
	}
	return ErrCannotPay
}

func (o *Order) Ship() error {
	if o.status == StatusPaid {
		o.status = StatusShipped
		return nil
	}
	return ErrCannotShip
}

func (o *Order) Cancel() error {
	if o.status == StatusShipped {
		return ErrCannotCancel
	}
	o.status = StatusCancelled
	return nil
}

func (o *Order) Total() Price {
	if len(o.items) == 0 {
		return 0
	}
	var total float64 = 0
	for _, item := range o.items {
		total += float64(item.price)
	}
	price, _ := NewPrice(total)
	return price
}

func (o *Order) Clone() *Order {
	copyOrder := *o

	itemsCopy := make([]*OrderItem, len(o.items))

	for i, item := range o.items {
		itemCopy := *item
		itemsCopy[i] = &itemCopy
	}

	copyOrder.items = itemsCopy

	return &copyOrder
}

func (o *Order) ID() OrderID {
	return o.id
}

func (o *Order) CustomerID() CustomerID {
	return o.customerID
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) Items() []*OrderItem {
	itemsCopy := make([]*OrderItem, len(o.items))

	for i, item := range o.items {
		itemCopy := *item
		itemsCopy[i] = &itemCopy
	}
	return itemsCopy
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}
