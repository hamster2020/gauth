package gauth

import "net/http"

type Authenticator interface {
	SetEmail(email string)
	SetPassword(password string)
	Authenticate() error
	DoAuthenticatedRequest(req *http.Request, respBody interface{}) error
}
