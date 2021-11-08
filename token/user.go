package token

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hamster2020/gauth"
)

type UserClaims struct {
	jwt.StandardClaims
	Email string      `json:"email"`
	Roles gauth.Roles `json:"roles"`
}

func (t token) NewUserToken(email string, roles gauth.Roles) (string, error) {
	stdClaims := newStandardClaims(email, time.Minute*time.Duration(t.accessTokenExpMinutes))
	userClaims := UserClaims{
		StandardClaims: stdClaims,
		Email:          email,
		Roles:          roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, userClaims)
	tokenStr, err := token.SignedString(t.accessTokenSigningKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (t token) VerifyUserToken(token string) (gauth.User, error) {
	var userClaims UserClaims
	parser := jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS512.Alg()}}
	if _, err := parser.ParseWithClaims(token, &userClaims, func(*jwt.Token) (interface{}, error) {
		return &t.accessTokenSigningKey.PublicKey, nil
	}); err != nil {
		return gauth.User{}, err
	}

	if err := validateStandardClaims(userClaims.StandardClaims); err != nil {
		return gauth.User{}, err
	}

	if userClaims.Email != userClaims.Subject {
		return gauth.User{}, fmt.Errorf("user claims subject %s does not match email %s", userClaims.Subject, userClaims.Email)
	}

	if err := userClaims.Roles.Validate(); err != nil {
		return gauth.User{}, err
	}

	user := gauth.User{
		Email: userClaims.Email,
		Roles: userClaims.Roles,
	}

	return user, nil
}
