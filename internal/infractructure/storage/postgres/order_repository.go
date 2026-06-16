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

type rawOrder struct {
	id         int64
	customerID int64
	status     string
	createdAt  time.Time
	nextItemID int64
	version    int64
}

type rawOrderItem struct {
	orderID   int64
	itemID    int64
	productID int64
	name      string
	price     int64
	quantity  int64
}

type persistedOrderItem struct {
	itemID    int64
	productID int64
	name      string
	price     int64
	quantity  int64
}

func restoreOrderItem(raw rawOrderItem) (*domain_order.OrderItem, error) {
	itemID, err := domain_order.NewItemID(raw.itemID)
	if err != nil {
		return nil, err
	}

	productID, err := domain_order.NewProductID(raw.productID)
	if err != nil {
		return nil, err
	}

	name, err := domain_order.NewProductName(raw.name)
	if err != nil {
		return nil, err
	}

	price, err := domain_order.NewPrice(raw.price)
	if err != nil {
		return nil, err
	}

	quantity, err := domain_order.NewQuantity(raw.quantity)
	if err != nil {
		return nil, err
	}

	return domain_order.NewOrderItem(itemID, productID, name, price, quantity), nil
}

func restoreOrder(raw rawOrder, items []*domain_order.OrderItem) (*domain_order.Order, error) {
	orderID, err := domain_order.NewOrderID(raw.id)
	if err != nil {
		return nil, err
	}

	customerID, err := domain_order.NewCustomerID(raw.customerID)
	if err != nil {
		return nil, err
	}

	status, err := domain_order.NewOrderStatus(raw.status)
	if err != nil {
		return nil, err
	}

	nextItemID, err := domain_order.NewItemID(raw.nextItemID)
	if err != nil {
		return nil, err
	}

	version, err := domain_order.NewOrderVersion(raw.version)
	if err != nil {
		return nil, err
	}

	return domain_order.RestoreOrder(
		orderID,
		customerID,
		status,
		raw.createdAt,
		nextItemID,
		version,
		items,
	), nil
}

func (r *PostgresOrderRepository) Create(ctx context.Context, customerID domain_order.CustomerID) (domain_order.OrderID, error) {
	status, err := domain_order.NewOrderStatus("created")
	if err != nil {
		return 0, err
	}
	query := `
		INSERT INTO orders (customer_id, status, created_at, next_item_id, version)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
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

	var raw rawOrder

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&raw.id,
		&raw.customerID,
		&raw.status,
		&raw.createdAt,
		&raw.nextItemID,
		&raw.version,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get order: %w", err)
	}

	query = `
		SELECT order_id, item_id, product_id, name, price_cents, quantity
		FROM order_items
		WHERE order_id = $1
		ORDER BY item_id
	`

	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("get order items: %w", err)
	}
	defer rows.Close()

	items := make([]*domain_order.OrderItem, 0)

	for rows.Next() {
		var rawItem rawOrderItem

		if err := rows.Scan(
			&rawItem.orderID,
			&rawItem.itemID,
			&rawItem.productID,
			&rawItem.name,
			&rawItem.price,
			&rawItem.quantity,
		); err != nil {
			return nil, fmt.Errorf("get order item: %w", err)
		}

		item, err := restoreOrderItem(rawItem)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order items: %w", err)
	}

	order, err := restoreOrder(raw, items)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *PostgresOrderRepository) List(ctx context.Context, params domain_order.OrderListParams) ([]*domain_order.Order, error) {
	var rows pgx.Rows
	var err error
	if params.HasCursor() {
		query := `
		SELECT id, customer_id, status, created_at, next_item_id, version
		FROM orders
		WHERE customer_id = $1 AND id > $2
		ORDER BY id
		LIMIT $3
	`

		rows, err = r.pool.Query(ctx, query, params.CustomerID(), params.Cursor(), params.Limit())
		if err != nil {
			return nil, fmt.Errorf("get orders: %w", err)
		}
	} else {
		query := `
		SELECT id, customer_id, status, created_at, next_item_id, version
		FROM orders
		WHERE customer_id = $1
		ORDER BY id
		LIMIT $2
	`

		rows, err = r.pool.Query(ctx, query, params.CustomerID(), params.Limit())
		if err != nil {
			return nil, fmt.Errorf("get orders: %w", err)
		}
	}

	defer rows.Close()

	rawOrders := make([]rawOrder, 0)
	orderIDs := make([]int64, 0)

	for rows.Next() {
		var raw rawOrder

		if err := rows.Scan(
			&raw.id,
			&raw.customerID,
			&raw.status,
			&raw.createdAt,
			&raw.nextItemID,
			&raw.version,
		); err != nil {
			return nil, fmt.Errorf("scan order: %w", err)
		}

		rawOrders = append(rawOrders, raw)
		orderIDs = append(orderIDs, raw.id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate orders: %w", err)
	}

	if len(rawOrders) == 0 {
		return []*domain_order.Order{}, nil
	}

	itemsQuery := `
		SELECT order_id, item_id, product_id, name, price_cents, quantity
		FROM order_items
		WHERE order_id = ANY($1)
		ORDER BY order_id, item_id
	`

	itemRows, err := r.pool.Query(ctx, itemsQuery, orderIDs)
	if err != nil {
		return nil, fmt.Errorf("get all order items: %w", err)
	}
	defer itemRows.Close()

	itemsByOrderID := make(map[int64][]*domain_order.OrderItem, len(rawOrders))

	for itemRows.Next() {
		var rawItem rawOrderItem

		if err := itemRows.Scan(
			&rawItem.orderID,
			&rawItem.itemID,
			&rawItem.productID,
			&rawItem.name,
			&rawItem.price,
			&rawItem.quantity,
		); err != nil {
			return nil, fmt.Errorf("scan order item: %w", err)
		}

		item, err := restoreOrderItem(rawItem)
		if err != nil {
			return nil, err
		}

		itemsByOrderID[rawItem.orderID] = append(itemsByOrderID[rawItem.orderID], item)
	}

	if err := itemRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate order items: %w", err)
	}

	orders := make([]*domain_order.Order, 0, len(rawOrders))

	for _, raw := range rawOrders {
		items := itemsByOrderID[raw.id]
		if items == nil {
			items = make([]*domain_order.OrderItem, 0)
		}

		order, err := restoreOrder(raw, items)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *PostgresOrderRepository) Update(ctx context.Context, order *domain_order.Order) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	query := `
		UPDATE orders
		SET status = $1, next_item_id = $2, version = version + 1
		WHERE id = $3 AND version = $4
	`
	tag, err := tx.Exec(ctx, query, string(order.Status()), int64(order.NextItemID()), int64(order.ID()), int64(order.Version()))
	if err != nil {
		return fmt.Errorf("update order row: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain_order.ErrConcurrentUpdate
	}

	dbItemsByID := make(map[int64]persistedOrderItem)

	query = `
		SELECT item_id, product_id, name, price_cents, quantity
		FROM order_items
		WHERE order_id = $1
	`
	rows, err := tx.Query(ctx, query, int64(order.ID()))
	if err != nil {
		return fmt.Errorf("get order items rows: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var item persistedOrderItem
		if err = rows.Scan(&item.itemID, &item.productID, &item.name, &item.price, &item.quantity); err != nil {
			return fmt.Errorf("parse order items row: %w", err)
		}
		dbItemsByID[item.itemID] = item
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate order items: %w", err)
	}

	actualItemsByID := make(map[int64]*domain_order.OrderItem)

	for _, item := range order.Items() {
		actualItemsByID[int64(item.ID())] = item
	}

	toDelete := make([]int64, 0)
	for itemID := range dbItemsByID {
		if _, exists := actualItemsByID[itemID]; !exists {
			toDelete = append(toDelete, itemID)
		}
	}
	toInsert := make([]*domain_order.OrderItem, 0)
	toUpdate := make([]*domain_order.OrderItem, 0)
	for itemID, item := range actualItemsByID {
		dbItem, exists := dbItemsByID[itemID]
		if !exists {
			toInsert = append(toInsert, item)
			continue
		}
		if !isSameOrderItem(dbItem, item) {
			toUpdate = append(toUpdate, item)
		}
	}
	deleteQuery := `
		DELETE FROM order_items
		WHERE order_id = $1 AND item_id = $2
	`
	for _, itemID := range toDelete {
		tag, err := tx.Exec(ctx, deleteQuery, int64(order.ID()), itemID)
		if err != nil {
			return fmt.Errorf("delete order item %d: %w", itemID, err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("delete order item %d: no rows affected", itemID)
		}
	}

	updateItemQuery := `
		UPDATE order_items
		SET product_id = $1,
		name = $2,
		price_cents = $3,
		quantity = $4
		WHERE order_id = $5 AND item_id = $6
	`
	for _, item := range toUpdate {
		tag, err := tx.Exec(ctx, updateItemQuery,
			int64(item.ProductID()),
			string(item.Name()),
			int64(item.Price()),
			int64(item.Quantity()),
			int64(order.ID()),
			int64(item.ID()),
		)
		if err != nil {
			return fmt.Errorf("update order item %d: %w", item.ID(), err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("update order item %d: no rows affected", item.ID())
		}
	}
	insertQuery := `
	INSERT INTO order_items (
		order_id,
		item_id,
		product_id,
		name,
		price_cents,
		quantity
		) VALUES ($1, $2, $3, $4, $5, $6)
	`
	for _, item := range toInsert {
		tag, err := tx.Exec(ctx, insertQuery,
			int64(order.ID()),
			int64(item.ID()),
			int64(item.ProductID()),
			string(item.Name()),
			int64(item.Price()),
			int64(item.Quantity()),
		)
		if err != nil {
			return fmt.Errorf("insert order item %d: %w", item.ID(), err)
		}
		if tag.RowsAffected() == 0 {
			return fmt.Errorf("insert order item %d: no rows affected", item.ID())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func isSameOrderItem(db persistedOrderItem, actual *domain_order.OrderItem) bool {
	if db.name != string(actual.Name()) {
		return false
	}
	if db.price != int64(actual.Price()) {
		return false
	}
	if db.productID != int64(actual.ProductID()) {
		return false
	}
	if db.quantity != int64(actual.Quantity()) {
		return false
	}
	return true
}
