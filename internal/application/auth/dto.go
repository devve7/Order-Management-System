package auth

import domain_user "Order-Management-System/internal/domain/user"

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	ID domain_user.UserID
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken AccessToken
}
