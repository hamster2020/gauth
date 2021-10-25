package logic

import (
	"github.com/hamster2020/gauth"
	"golang.org/x/crypto/bcrypt"
)

func checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (l logic) Authenticate(c gauth.Credentials) (bool, error) {
	return authenticate(l.ds, c, checkPassword)
}

func authenticate(
	ds gauth.Datastore,
	c gauth.Credentials,
	checkPasswordFunc func(password, hash string) bool,
) (bool, error) {
	user, err := ds.UserByEmail(c.Email)
	if err != nil {
		return false, err
	}

	return checkPasswordFunc(c.Password, user.PasswordHash), nil
}
