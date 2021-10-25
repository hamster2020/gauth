package mapstore

import "errors"

type mapStore map[string]string

var notFoundErr = errors.New("not found")

func NewMapStore() mapStore {
	return make(map[string]string)
}

func (m mapStore) set(key, value string) {
	m[key] = value
}

func (m mapStore) get(key string) (string, error) {
	value, found := m[key]
	if !found {
		return "", notFoundErr
	}

	return value, nil
}

func (m mapStore) delete(key string) {
	delete(m, key)
}
