package user

type User struct {
	id           UserID
	email        Email
	passwordHash PasswordHash
	role         Role
}

func NewUser(email Email, passwordHash PasswordHash, role Role) *User {
	return &User{
		id:           0,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
	}
}

func RestoreUser(id UserID, email Email, passwordHash PasswordHash, role Role) *User {
	return &User{
		id:           id,
		email:        email,
		passwordHash: passwordHash,
		role:         role,
	}
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

func (u *User) Role() Role {
	return u.role
}
