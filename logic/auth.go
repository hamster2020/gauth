package logic

import (
	"errors"

	"github.com/hamster2020/gauth"
	"golang.org/x/crypto/bcrypt"
)

func checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (l logic) Authenticate(c gauth.Credentials) (string, error) {
	return authenticate(l.ds, l.token, c, checkPassword)
}

func authenticate(
	ds gauth.Datastore,
	token gauth.Token,
	c gauth.Credentials,
	checkPasswordFunc func(password, hash string) bool,
) (string, error) {
	user, err := ds.UserByEmail(c.Email)
	if err != nil {
		return "", err
	}

	if pass := checkPasswordFunc(c.Password, user.PasswordHash); !pass {
		return "", errors.New("unauthorized")
	}

	tokenStr, err := token.NewUserToken(user.Email, user.Roles)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
