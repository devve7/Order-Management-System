// Package postgres
package postgres

import (
	domain_order "Order-Management-System/internal/domain/order"
	"context"
	"fmt"
	"time"

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

	now := time.Now()

	var id int64
	err = r.pool.QueryRow(
		ctx,
		`
		INSERT INTO orders (customer_id, status, created_at, next_item_id, version)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
		`,
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
	return nil, nil
}

func (r *PostgresOrderRepository) Update(ctx context.Context, order *domain_order.Order) error {
	return nil
}

func (r *PostgresOrderRepository) GetAll(ctx context.Context) ([]*domain_order.Order, error) {
	return nil, nil
}
