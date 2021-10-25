package gauth

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLogic interface {
	// General Use
	Authenticate(c Credentials) (bool, error)
}
