package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (UserID, error)
	GetByEmail(ctx context.Context, email Email) (*User, error)
	GetByID(ctx context.Context, id UserID) (*User, error)
}
