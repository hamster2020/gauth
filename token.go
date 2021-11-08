package gauth

import "crypto/rsa"

type Token interface {
	PublicKey() rsa.PublicKey
	NewUserToken(email string, roles Roles) (string, error)
	VerifyUserToken(tokenStr string) (User, error)
}
