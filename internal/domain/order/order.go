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
	nextItemID ItemID
	version    OrderVersion
}

func NewOrder(id OrderID, customerID CustomerID, status OrderStatus, time time.Time) *Order {
	return &Order{
		id:         id,
		customerID: customerID,
		items:      make([]*OrderItem, 0),
		status:     status,
		createdAt:  time,
		nextItemID: 1,
		version:    1,
	}
}

func RestoreOrder(id OrderID, customerID CustomerID, status OrderStatus, time time.Time, nextItemID ItemID, version OrderVersion, items []*OrderItem) *Order {
	return &Order{
		id:         id,
		customerID: customerID,
		items:      items,
		status:     status,
		createdAt:  time,
		nextItemID: nextItemID,
		version:    version,
	}
}

func (o *Order) AddItem(productID ProductID, name ProductName, price Price, quantity Quantity) error {
	if o.status == StatusCancelled {
		return ErrOrderCancelled
	}
	if o.status == StatusShipped {
		return ErrOrderShipped
	}

	itemID := o.getNextItemID()
	item := NewOrderItem(itemID, productID, name, price, quantity)
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
	if o.status == StatusShipped || o.status == StatusCancelled {
		return ErrCannotCancel
	}
	o.status = StatusCancelled
	return nil
}

func (o *Order) Total() Price {
	if len(o.items) == 0 {
		return 0
	}
	var total int64 = 0
	for _, item := range o.items {
		total += int64(item.Price()) * int64(item.Quantity())
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

func (o *Order) getNextItemID() ItemID {
	id := o.nextItemID
	o.nextItemID++
	return id
}

func (o *Order) Version() OrderVersion {
	return o.version
}

func (o *Order) WithNextVersion() *Order {
	clone := o.Clone()
	clone.version++
	return clone
}
