package auth

import (
	domain_user "Order-Management-System/internal/domain/user"
	"context"
	"errors"
)

type AuthService struct {
	repo           domain_user.Repository
	passwordHasher PasswordHasher
	tokenManager   TokenManager
}

func NewAuthService(repo domain_user.Repository, passwordHasher PasswordHasher, tokenManager TokenManager) *AuthService {
	return &AuthService{
		repo:           repo,
		passwordHasher: passwordHasher,
		tokenManager:   tokenManager,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (RegisterOutput, error) {
	email, err := domain_user.NewEmail(input.Email)
	if err != nil {
		return RegisterOutput{}, err
	}
	passwordHash, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return RegisterOutput{}, err
	}
	user := domain_user.NewUser(email, passwordHash, domain_user.RoleUser)
	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return RegisterOutput{}, err
	}
	return RegisterOutput{
		ID: userID,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (LoginOutput, error) {
	email, err := domain_user.NewEmail(input.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain_user.ErrUserNotFound) {
			return LoginOutput{}, ErrInvalidCredentials
		}
		return LoginOutput{}, err
	}
	hash := user.PasswordHash()
	if err = s.passwordHasher.Compare(input.Password, hash); err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}
	actor := NewActor(user.ID(), user.Role())
	accessToken, err := s.tokenManager.CreateAccessToken(actor)
	if err != nil {
		return LoginOutput{}, err
	}
	return LoginOutput{
		AccessToken: accessToken,
	}, nil
}
