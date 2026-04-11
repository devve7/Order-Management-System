// Package postgres
package postgres

import (
	domain_order "Order-Management-System/internal/domain/order"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		pool: pool,
	}
}

func (r *PostgresOrderRepository) Create(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
	status, err := domain_order.NewOrderStatus("created")
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO orders (customer_id, status, created_at, next_item_id, version)
		VALUES ($1, $2, $3, $4, $5)
		RETURING id
	`

	now := time.Now()

	var id int64
	err = r.pool.QueryRow(
		ctx,
		query,
		int64(customerID),
		string(status),
		now,
		int64(1), // next_item_id
		int64(1), // version
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert order: %w", err)
	}

	orderID, err := domain_order.NewOrderID(id)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func (r *PostgresOrderRepository) Get(ctx context.Context, id domain_order.OrderID) (*domain_order.Order, error) {
	query := `
		SELECT id, customer_id, status, created_at, next_item_id, version
		FROM orders
		WHERE id = $1
	`
	var rawOrderID int64
	var rawCustomerID int64
	var rawStatus string
	var rawCreatedAt time.Time
	var rawNextItemID int64
	var rawVersion int64

	err := r.pool.QueryRow(ctx, query, id).Scan(&rawOrderID, &rawCustomerID, &rawStatus, &rawCreatedAt, &rawNextItemID, &rawVersion)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order: %w", err)
	}

	orderID, err := domain_order.NewOrderID(rawOrderID)
	if err != nil {
		return nil, err
	}
	customerID, err := domain_order.NewCustomerID(rawCustomerID)
	if err != nil {
		return nil, err
	}
	status, err := domain_order.NewOrderStatus(rawStatus)
	if err != nil {
		return nil, err
	}
	nextItemID, err := domain_order.NewItemID(rawNextItemID)
	if err != nil {
		return nil, err
	}
	version, err := domain_order.NewOrderVersion(rawVersion)
	if err != nil {
		return nil, err
	}

	query = `
		SELECT item_id, product_id, name, price, quantity
		FROM order_items 
		WHERE order_id = $1
	`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()
	items := make([]*domain_order.OrderItem, 0)
	for rows.Next() {
		var rawItemID int64
		var rawProductID int64
		var rawName string
		var rawPrice int64
		var rawQuantity int64
		if err := rows.Scan(&rawItemID, &rawProductID, &rawName, &rawPrice, &rawQuantity); err != nil {
			return nil, fmt.Errorf("get order item: %w", err)
		}
		itemID, err := domain_order.NewItemID(rawItemID)
		if err != nil {
			return nil, err
		}
		productID, err := domain_order.NewProductID(rawProductID)
		if err != nil {
			return nil, err
		}
		name, err := domain_order.NewProductName(rawName)
		if err != nil {
			return nil, err
		}
		price, err := domain_order.NewPrice(rawPrice)
		if err != nil {
			return nil, err
		}
		quantity, err := domain_order.NewQuantity(rawQuantity)
		if err != nil {
			return nil, err
		}
		orderItem := domain_order.NewOrderItem(itemID, productID, name, price, quantity)
		items = append(items, orderItem)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order items: %w", err)
	}

	order := domain_order.RestoreOrder(orderID, customerID, status, rawCreatedAt, nextItemID, version, items)

	return order, nil
}

func (r *PostgresOrderRepository) Update(ctx context.Context, order *domain_order.Order) error {
	return nil
}

func (r *PostgresOrderRepository) GetAll(ctx context.Context) ([]*domain_order.Order, error) {
	return nil, nil

}
