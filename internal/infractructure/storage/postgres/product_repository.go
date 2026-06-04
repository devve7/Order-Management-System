package postgres

import (
	domain_product "Order-Management-System/internal/domain/product"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type rawProduct struct {
	id     int64
	name   string
	price  int64
	stock  int64
	active bool
}

type PostgresProductRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresProductRepository(pool *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{
		pool: pool,
	}
}

func (r *PostgresProductRepository) Create(ctx context.Context, product *domain_product.Product) (domain_product.ProductID, error) {
	if product.HasID() {
		return 0, errors.New("product already has id")
	}
	query := `
		INSERT INTO products(name, price_cents, stock, active)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var id int64
	err := r.pool.QueryRow(
		ctx,
		query,
		string(product.Name()),
		int64(product.Price()),
		int64(product.Stock()),
		product.Active(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert product: %w", err)
	}
	productID, err := domain_product.NewProductID(id)
	if err != nil {
		return 0, err
	}
	return productID, nil
}

func restoreProduct(rawProduct rawProduct) (*domain_product.Product, error) {
	id, err := domain_product.NewProductID(rawProduct.id)
	if err != nil {
		return nil, err
	}
	name, err := domain_product.NewProductName(rawProduct.name)
	if err != nil {
		return nil, err
	}
	price, err := domain_product.NewPrice(rawProduct.price)
	if err != nil {
		return nil, err
	}
	stock, err := domain_product.NewStock(rawProduct.stock)
	if err != nil {
		return nil, err
	}
	product := domain_product.RestoreProduct(id, name, price, stock, rawProduct.active)
	return product, nil
}

func (r *PostgresProductRepository) Get(ctx context.Context, id domain_product.ProductID) (*domain_product.Product, error) {
	query := `
		SELECT id, name, price_cents, stock, active 
		FROM products
		WHERE id = $1
	`
	var rawProduct rawProduct
	err := r.pool.QueryRow(ctx, query, int64(id)).Scan(
		&rawProduct.id,
		&rawProduct.name,
		&rawProduct.price,
		&rawProduct.stock,
		&rawProduct.active,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain_product.ErrProductNotFound
		}
		return nil, fmt.Errorf("get product: %w", err)
	}
	return restoreProduct(rawProduct)
}

func (r *PostgresProductRepository) List(ctx context.Context, params domain_product.ProductListParams) ([]*domain_product.Product, error) {
	products := make([]*domain_product.Product, 0)
	var query string
	var rows pgx.Rows
	var err error
	if params.HasCursor() {
		query = `
			SELECT id, name, price_cents, stock, active 
			FROM products
			WHERE id > $1
			ORDER BY id
			LIMIT $2
		`
		rows, err = r.pool.Query(ctx, query, params.GetCursor(), params.GetLimit())
		if err != nil {
			return nil, fmt.Errorf("get products: %w", err)
		}
	} else {
		query = `
			SELECT id, name, price_cents, stock, active 
			FROM products
			ORDER BY id
			LIMIT $1
		`
		rows, err = r.pool.Query(ctx, query, params.GetLimit())
		if err != nil {
			return nil, fmt.Errorf("get products: %w", err)
		}
	}
	defer rows.Close()
	for rows.Next() {
		var rawProduct rawProduct
		if err := rows.Scan(
			&rawProduct.id,
			&rawProduct.name,
			&rawProduct.price,
			&rawProduct.stock,
			&rawProduct.active,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}
		product, err := restoreProduct(rawProduct)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}
	return products, nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, product *domain_product.Product) error {
	query := `
		UPDATE products
		SET name = $1, price_cents = $2, stock = $3, active = $4
		WHERE id = $5
	`
	tag, err := r.pool.Exec(ctx,
		query,
		string(product.Name()),
		int64(product.Price()),
		int64(product.Stock()),
		product.Active(),
		int64(product.ID()),
	)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain_product.ErrProductNotFound
	}
	return nil
}
