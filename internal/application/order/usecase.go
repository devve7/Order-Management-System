package order

import (
	"Order-Management-System/internal/domain/order"
)

type UseCase struct {
	repo order.Repository
}

func NewUseCase(repo order.Repository) *UseCase {
	return &UseCase{
		repo: repo,
	}
}

func (u *UseCase) CreateOrder(customerID order.CustomerID) (order.OrderID, error) {
	orderID, err := u.repo.NextID()
	if err != nil {
		return 0, err
	}
	order := order.NewOrder(orderID, customerID)
	err = u.repo.Save(order)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *UseCase) AddItem(orderID order.OrderID, item *order.OrderItem) error {
	order, err := u.repo.Get(orderID)
	if err != nil {
		return err
	}
	err = order.AddItem(item)
	if err != nil {
		return err
	}
	err = u.repo.Update(order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) RemoveItem(orderID order.OrderID, itemID order.ItemID) error {
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

// PayOrder(orderID string) error

func (u *UseCase) PayOrder(orderID order.OrderID) error {
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

func (u *UseCase) ShipOrder(orderID order.OrderID) error {
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

func (u *UseCase) CancelOrder(orderID order.OrderID) error {
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
