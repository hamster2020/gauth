package mapstore

import (
	"errors"

	"github.com/hamster2020/gauth"
)

var userExistsErr = errors.New("user with email already exists")

func newCred(email, hash string) gauth.User {
	return gauth.User{Email: email, PasswordHash: hash}
}

func (m mapStore) CreateUser(user gauth.User) error {
	if hash, err := m.get(user.Email); hash != "" && err == nil {
		return userExistsErr
	}

	m.set(user.Email, user.PasswordHash)
	return nil
}

func (m mapStore) UserByEmail(email string) (gauth.User, error) {
	hash, err := m.get(email)
	if err != nil {
		return gauth.User{}, err
	}

	return newCred(email, hash), nil
}

func (m mapStore) UpdateUser(email string, user gauth.User) error {
	_, err := m.get(email)
	if err != nil {
		return err
	}

	if email == user.Email {
		m.set(email, user.PasswordHash)
		return nil
	}

	m.delete(email)
	m.set(user.Email, user.PasswordHash)
	return nil
}

func (m mapStore) DeleteUser(email string) error {
	m.delete(email)
	return nil
}

func (m mapStore) Users() ([]gauth.User, error) {
	ret := make([]gauth.User, len(m))
	i := 0
	for email, hash := range m {
		ret[i] = gauth.User{Email: email, PasswordHash: hash}
		i++
	}

	return ret, nil
}
