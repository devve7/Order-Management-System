// Package order ...
package order

import (
	"errors"
	"sync"
	"time"
)

type Order struct {
	id         OrderID
	customerID CustomerID
	items      []*OrderItem
	status     OrderStatus
	createdAt  time.Time
	mtx        sync.RWMutex
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
	o.mtx.Lock()
	defer o.mtx.Unlock()
	if item == nil {
		return errors.New("item cannot be nil")
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
	o.mtx.Lock()
	defer o.mtx.Unlock()
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
		if item != nil && item.id == id {
			o.items = append(o.items[:index], o.items[index+1:]...)
			return nil
		}
	}
	return ErrOrderItemNotFound
}

func (o *Order) Pay() error {
	o.mtx.Lock()
	defer o.mtx.Unlock()
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
	o.mtx.Lock()
	defer o.mtx.Unlock()
	if o.status == StatusPaid {
		o.status = StatusShipped
		return nil
	}
	return ErrCannotShip
}

func (o *Order) Cancel() error {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	if o.status == StatusShipped {
		return ErrCannotCancel
	}
	o.status = StatusCancelled
	return nil
}

func (o *Order) Total() (Price, error) {
	o.mtx.RLock()
	defer o.mtx.RUnlock()
	if o.items == nil {
		return 0, ErrOrderEmpty
	}
	if len(o.items) == 0 {
		return 0, ErrOrderEmpty
	}
	var total float64 = 0
	for _, item := range o.items {
		total += float64(item.price)
	}
	return NewPrice(total)
}
