package mapstore

import (
	"errors"

	"github.com/hamster2020/gauth"
)

var sessionExistsErr = errors.New("session with cookie already exists")

func (m mapStore) CreateSession(session gauth.Session) error {
	_, err := m.getSession(session.Cookie)
	if err == nil {
		return sessionExistsErr
	}
	if err != notFoundErr {
		return err
	}

	m.setSession(session.Cookie, session)
	return nil
}

func (m mapStore) SessionByCookie(cookie string) (gauth.Session, error) {
	session, err := m.getSession(cookie)
	if err != nil {
		return gauth.Session{}, err
	}

	return session, nil
}

func (m mapStore) DeleteSession(cookie string) error {
	delete(m.sessions, cookie)
	return nil
}
