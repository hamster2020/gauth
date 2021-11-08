package gauth

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	Roles        Roles  `json:"roles"`
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
