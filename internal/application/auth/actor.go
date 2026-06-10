package auth

import domain_user "Order-Management-System/internal/domain/user"

type Actor struct {
	id   domain_user.UserID
	role domain_user.Role
}

func NewActor(id domain_user.UserID, role domain_user.Role) Actor {
	return Actor{
		id:   id,
		role: role,
	}
}

func (a *Actor) ID() domain_user.UserID {
	return a.id
}

func (a *Actor) Role() domain_user.Role {
	return a.role
}
