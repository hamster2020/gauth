package gauth

import "golang.org/x/crypto/bcrypt"

type User struct {
	Email        string `db:"email" json:"email"`
	PasswordHash string `db:"password_hash" json:"-"`
	Roles        Roles  `db:"roles" json:"roles"`
}

func (u User) Info() string {
	return u.Email
}

func (u User) HasRole(role Roles) bool {
	return u.Roles.HasRole(role)
}

func (u User) HasAtLeastOneRole(roles Roles) bool {
	return u.Roles.HasAtLeastOneRole(roles)
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Roles    Roles  `json:"roles"`
}

type UserLogic interface {
	// CRUD
	CreateUser(req UserRequest) error
	UserByEmail(email string) (User, error)
	UpdateUser(oldEmail string, req UserRequest) error
	DeleteUser(email string) error

	// List
	Users() ([]User, error)
}

type UserStore interface {
	// CRUD
	CreateUser(u User) error
	UserByEmail(email string) (User, error)
	UpdateUser(email string, u User) error
	DeleteUser(email string) error

	// List
	Users() ([]User, error)
}
