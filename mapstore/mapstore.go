package mapstore

import (
	"errors"

	"github.com/hamster2020/gauth"
)

type mapStore struct {
	users    map[string]gauth.User
	sessions map[string]gauth.Session
}

var notFoundErr = errors.New("not found")

func NewMapStore() mapStore {
	return mapStore{
		users:    make(map[string]gauth.User),
		sessions: make(map[string]gauth.Session),
	}
}

func (m mapStore) setUser(key string, value gauth.User) {
	m.users[key] = value
}

func (m mapStore) getUser(key string) (gauth.User, error) {
	value, found := m.users[key]
	if !found {
		return gauth.User{}, notFoundErr
	}

	return value, nil
}

func (m mapStore) deleteUser(key string) {
	delete(m.users, key)
}

func (m mapStore) setSession(key string, value gauth.Session) {
	m.sessions[key] = value
}

func (m mapStore) getSession(key string) (gauth.Session, error) {
	value, found := m.sessions[key]
	if !found {
		return gauth.Session{}, notFoundErr
	}

	return value, nil
}

func (m mapStore) deleteSession(key string) {
	delete(m.sessions, key)
}
