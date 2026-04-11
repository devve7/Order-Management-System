// Package postgres
package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type PostgresOrderRepository struct {
	pool *pgxpool.Pool
}
