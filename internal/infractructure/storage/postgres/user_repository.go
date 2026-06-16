package postgres

import (
	user_domain "Order-Management-System/internal/domain/user"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{
		pool: pool,
	}
}

type rawUser struct {
	id           int64
	email        string
	passwordHash string
	role         string
}

func restoreUser(rawUser rawUser) (*user_domain.User, error) {
	id, err := user_domain.NewUserID(rawUser.id)
	if err != nil {
		return nil, err
	}
	email, err := user_domain.NewEmail(rawUser.email)
	if err != nil {
		return nil, err
	}
	passwordHash, err := user_domain.NewPasswordHash(rawUser.passwordHash)
	if err != nil {
		return nil, err
	}
	role, err := user_domain.NewRole(rawUser.role)
	if err != nil {
		return nil, err
	}
	user := user_domain.RestoreUser(id, email, passwordHash, role)
	return user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *user_domain.User) (user_domain.UserID, error) {
	if user.HasID() {
		return 0, user_domain.ErrUserAlreadyHasID
	}
	query := `
		INSERT INTO users(email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int64
	err := r.pool.QueryRow(
		ctx,
		query,
		string(user.Email()),
		string(user.PasswordHash()),
		string(user.Role()),
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, user_domain.ErrUserAlreadyExists
		}

		return 0, fmt.Errorf("insert user: %w", err)
	}
	userID, err := user_domain.NewUserID(id)
	if err != nil {
		return 0, fmt.Errorf("parse user id from db: %w", err)
	}
	return userID, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email user_domain.Email) (*user_domain.User, error) {
	query := `
		SELECT id, email, password_hash, role
		FROM users
		WHERE email = $1
	`
	var raw rawUser
	err := r.pool.QueryRow(ctx, query, string(email)).Scan(
		&raw.id,
		&raw.email,
		&raw.passwordHash,
		&raw.role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return restoreUser(raw)
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id user_domain.UserID) (*user_domain.User, error) {
	query := `
		SELECT id, email, password_hash, role
		FROM users
		WHERE id = $1
	`
	var raw rawUser
	err := r.pool.QueryRow(ctx, query, int64(id)).Scan(
		&raw.id,
		&raw.email,
		&raw.passwordHash,
		&raw.role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user_domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return restoreUser(raw)
}
