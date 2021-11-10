package mapstore

import (
	"errors"
	"sort"

	"github.com/hamster2020/gauth"
)

var userExistsErr = errors.New("user with email already exists")

func newCred(email, hash string, roles gauth.Roles) gauth.User {
	return gauth.User{Email: email, PasswordHash: hash, Roles: roles}
}

func (m mapStore) CreateUser(user gauth.User) error {
	_, err := m.getUser(user.Email)
	if err == nil {
		return userExistsErr
	}
	if err != notFoundErr {
		return err
	}

	m.setUser(user.Email, user)
	return nil
}

func (m mapStore) UserByEmail(email string) (gauth.User, error) {
	user, err := m.getUser(email)
	if err != nil {
		return gauth.User{}, err
	}

	return user, nil
}

func (m mapStore) UpdateUser(email string, user gauth.User) error {
	_, err := m.getUser(email)
	if err != nil {
		return err
	}

	m.deleteUser(email)
	m.setUser(user.Email, user)
	return nil
}

func (m mapStore) DeleteUser(email string) error {
	m.deleteUser(email)
	return nil
}

func (m mapStore) Users() ([]gauth.User, error) {
	ret := make([]gauth.User, len(m.users))
	i := 0
	for _, user := range m.users {
		ret[i] = user
		i++
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Email < ret[j].Email })

	return ret, nil
}
