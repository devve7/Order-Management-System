// Package order ...
package order

import (
	do "Order-Management-System/internal/domain/order"
)

type UseCase struct {
	itemsFactory *do.OrderItemFactory
	repo         do.Repository
}

func NewUseCase(factory *do.OrderItemFactory, repo do.Repository) *UseCase {
	return &UseCase{
		itemsFactory: factory,
		repo:         repo,
	}
}

func (u *UseCase) CreateOrder(customerID do.CustomerID) (do.OrderID, error) {
	orderID, err := u.repo.NextID()
	if err != nil {
		return 0, err
	}
	order := do.NewOrder(orderID, customerID)
	err = u.repo.Save(order)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *UseCase) AddItem(orderID do.OrderID, name string) (do.ItemID, error) {
	item, err := u.itemsFactory.New(name)
	if err != nil {
		return do.ItemID(0), err
	}
	order, err := u.repo.Get(orderID)
	if err != nil {
		return do.ItemID(0), err
	}
	err = order.AddItem(item)
	if err != nil {
		return do.ItemID(0), err
	}
	err = u.repo.Update(order)
	if err != nil {
		return do.ItemID(0), err
	}
	return item.ID(), nil
}

func (u *UseCase) RemoveItem(orderID do.OrderID, itemID do.ItemID) error {
	order, err := u.repo.Get(orderID)
	if err != nil {
		return err
	}
	err = order.RemoveItem(itemID)
	if err != nil {
		return err
	}
	err = u.repo.Update(order)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) PayOrder(orderID do.OrderID) error {
	order, err := u.repo.Get(orderID)
	if err != nil {
		return err
	}
	err = order.Pay()
	if err != nil {
		return err
	}
	err = u.repo.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) ShipOrder(orderID do.OrderID) error {
	order, err := u.repo.Get(orderID)
	if err != nil {
		return err
	}
	err = order.Ship()
	if err != nil {
		return err
	}
	err = u.repo.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) CancelOrder(orderID do.OrderID) error {
	order, err := u.repo.Get(orderID)
	if err != nil {
		return err
	}
	err = order.Cancel()
	if err != nil {
		return err
	}
	err = u.repo.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) GetOrder(orderID do.OrderID) (OrderDTO, error) {
	order, err := u.repo.Get(orderID)
	if err != nil {
		return OrderDTO{}, err
	}
	orderDTO := ToOrderDTO(order)
	return orderDTO, nil
}

func (u *UseCase) GetOrders() ([]OrderDTO, error) {
	orders, err := u.repo.GetAll()
	if err != nil {
		return nil, err
	}

	ordersDTO := make([]OrderDTO, 0, len(orders))

	for _, order := range orders {
		ordersDTO = append(ordersDTO, ToOrderDTO(order))
	}

	return ordersDTO, nil
}
