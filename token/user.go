package token

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hamster2020/gauth"
)

type UserToken struct {
	jwt.StandardClaims
	Email string      `json:"email"`
	Roles gauth.Roles `json:"roles"`
}

func (t token) NewUserToken(email string, roles gauth.Roles) (string, error) {
	std := newStandardClaims(email, time.Minute*time.Duration(t.accessTokenExpMinutes))
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, std)
	tokenStr, err := token.SignedString(t.accessTokenSigningKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
