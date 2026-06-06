// Package order ...
package order

import (
	"Order-Management-System/internal/application/ports"
	domain_order "Order-Management-System/internal/domain/order"
	"context"
)

type UseCase struct {
	repo           domain_order.Repository
	productService ports.ProductForOrder
}

func NewUseCase(productService ports.ProductForOrder, repo domain_order.Repository) *UseCase {
	return &UseCase{
		repo:           repo,
		productService: productService,
	}
}

func (u *UseCase) CreateOrder(ctx context.Context, rawCustomerID int64) (domain_order.OrderID, error) {
	customerID, err := domain_order.NewCustomerID(rawCustomerID)
	if err != nil {
		return 0, err
	}
	orderID, err := u.repo.Create(ctx, customerID)
	if err != nil {
		return 0, err
	}
	return orderID, nil
}

func (u *UseCase) AddItem(ctx context.Context, rawOrderID int64, rawProductID int64, rawQuantity int64) error {
	orderID, err := domain_order.NewOrderID(rawOrderID)
	if err != nil {
		return err
	}
	productID, err := domain_order.NewProductID(rawProductID)
	if err != nil {
		return err
	}
	quantity, err := domain_order.NewQuantity(rawQuantity)
	if err != nil {
		return err
	}
	err = u.productService.EnsureAvailable(ctx, int64(productID), int64(quantity))
	if err != nil {
		return err
	}
	snapshot, err := u.productService.GetSnapshot(ctx, int64(productID))
	if err != nil {
		return err
	}
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return err
	}
	name, err := domain_order.NewProductName(snapshot.Name)
	if err != nil {
		return err
	}
	price, err := domain_order.NewPrice(snapshot.Price)
	if err != nil {
		return err
	}
	err = order.AddItem(productID, name, price, quantity)
	if err != nil {
		return err
	}
	err = u.repo.Update(ctx, order)
	if err != nil {
		return err
	}
	return nil
}

func (u *UseCase) RemoveItem(ctx context.Context, rawOrderID int64, rawItemID int64) error {
	orderID, err := domain_order.NewOrderID(rawOrderID)
	if err != nil {
		return err
	}
	itemID, err := domain_order.NewItemID(rawItemID)
	if err != nil {
		return err
	}
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

func (u *UseCase) PayOrder(ctx context.Context, rawOrderID int64) error {
	orderID, err := domain_order.NewOrderID(rawOrderID)
	if err != nil {
		return err
	}
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

func (u *UseCase) ShipOrder(ctx context.Context, rawOrderID int64) error {
	orderID, err := domain_order.NewOrderID(rawOrderID)
	if err != nil {
		return err
	}
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

func (u *UseCase) CancelOrder(ctx context.Context, rawOrderID int64) error {
	orderID, err := domain_order.NewOrderID(rawOrderID)
	if err != nil {
		return err
	}
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

func (u *UseCase) GetOrder(ctx context.Context, rawOrderID int64) (OrderDTO, error) {
	orderID, err := domain_order.NewOrderID(rawOrderID)
	if err != nil {
		return OrderDTO{}, err
	}
	order, err := u.repo.Get(ctx, orderID)
	if err != nil {
		return OrderDTO{}, err
	}
	orderDTO := ToOrderDTO(order)
	return orderDTO, nil
}

func (u *UseCase) GetOrders(ctx context.Context, paramsDTO OrderListParamsDTO) ([]OrderDTO, error) {
	params, err := domain_order.NewOrderListParams(paramsDTO.Cursor, paramsDTO.Limit)
	if err != nil {
		return nil, err
	}
	orders, err := u.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	ordersDTO := make([]OrderDTO, 0, len(orders))

	for _, order := range orders {
		ordersDTO = append(ordersDTO, ToOrderDTO(order))
	}

	return ordersDTO, nil
}
