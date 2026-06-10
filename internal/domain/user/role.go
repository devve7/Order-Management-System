package user

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func NewRole(role string) (Role, error) {
	switch Role(role) {
	case RoleUser, RoleAdmin:
		return Role(role), nil
	default:
		return "", ErrInvalidRole
	}
}

func (r Role) IsAdmin() bool {
	return r == RoleAdmin
}

func (r Role) IsUser() bool {
	return r == RoleUser
}
