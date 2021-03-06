// Code generated by mockit.
// DO NOT EDIT!

package mocks

import (
	"crypto/rsa"

	"github.com/hamster2020/gauth"
)

type MockToken struct {
	PublicKeyFunc       func() rsa.PublicKey
	NewUserTokenFunc    func(email string, roles gauth.Roles) (string, error)
	VerifyUserTokenFunc func(tokenStr string) (gauth.User, error)
}

func NewMockToken() *MockToken {
	return &MockToken{}
}

func (token *MockToken) PublicKey() rsa.PublicKey {
	return token.PublicKeyFunc()
}

func (token *MockToken) NewUserToken(email string, roles gauth.Roles) (string, error) {
	return token.NewUserTokenFunc(email, roles)
}

func (token *MockToken) VerifyUserToken(tokenStr string) (gauth.User, error) {
	return token.VerifyUserTokenFunc(tokenStr)
}
