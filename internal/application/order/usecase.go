// Package order ...
package order

import (
	product "Order-Management-System/internal/application/product"
	domain_order "Order-Management-System/internal/domain/order"
	"context"
)

type UseCase struct {
	repo           domain_order.Repository
	productService product.ProductService
}

func NewUseCase(productService product.ProductService, repo domain_order.Repository) *UseCase {
	return &UseCase{
		repo:           repo,
		productService: productService,
	}
}

func (u *UseCase) CreateOrder(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
	orderID, err := u.repo.Create(ctx, customerID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *UseCase) AddItem(ctx context.Context, orderID domain_order.OrderID, productID domain_order.ProductID, quantity domain_order.Quantity) error {
	err := u.productService.EnsureAvailable(ctx, productID, quantity)
	if err != nil {
		return err
	}
	snapshot, err := u.productService.GetSnapshot(ctx, productID)
	if err != nil {
		return err
	}
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return err
	}
	err = order.AddItem(productID, snapshot.Name, snapshot.Price, quantity)
	if err != nil {
		return err
	}
	err = u.repo.Update(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) RemoveItem(ctx context.Context, orderID domain_order.OrderID, itemID domain_order.ItemID) error {
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return err
	}
	err = order.RemoveItem(itemID)
	if err != nil {
		return err
	}
	err = u.repo.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}

func (u *UseCase) PayOrder(ctx context.Context, orderID domain_order.OrderID) error {
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return err
	}
	err = order.Pay()
	if err != nil {
		return err
	}
	err = u.repo.Update(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) ShipOrder(ctx context.Context, orderID domain_order.OrderID) error {
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return err
	}
	err = order.Ship()
	if err != nil {
		return err
	}
	err = u.repo.Update(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) CancelOrder(ctx context.Context, orderID domain_order.OrderID) error {
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return err
	}
	err = order.Cancel()
	if err != nil {
		return err
	}
	err = u.repo.Update(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) GetOrder(ctx context.Context, orderID domain_order.OrderID) (OrderDTO, error) {
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return OrderDTO{}, err
	}
	orderDTO := ToOrderDTO(order)
	return orderDTO, nil
}

func (u *UseCase) GetOrders(ctx context.Context) ([]OrderDTO, error) {
	orders, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	ordersDTO := make([]OrderDTO, 0, len(orders))

	for _, order := range orders {
		ordersDTO = append(ordersDTO, ToOrderDTO(order))
	}

	return ordersDTO, nil
}
