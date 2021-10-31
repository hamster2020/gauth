package mapstore

import (
	"errors"

	"github.com/hamster2020/gauth"
)

type mapStore map[string]gauth.User

var notFoundErr = errors.New("not found")

func NewMapStore() mapStore {
	return make(map[string]gauth.User)
}

func (m mapStore) set(key string, value gauth.User) {
	m[key] = value
}

func (m mapStore) get(key string) (gauth.User, error) {
	value, found := m[key]
	if !found {
		return gauth.User{}, notFoundErr
	}

	return value, nil
}

func (m mapStore) delete(key string) {
	delete(m, key)
}
