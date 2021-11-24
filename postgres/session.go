package postgres

import (
	"database/sql"

	"github.com/hamster2020/gauth"
)

func (db DB) CreateSession(s gauth.Session) error {
	return createSession(db, s)
}

func (db DB) SessionByCookie(cookie string) (gauth.Session, error) {
	return sessionByCookie(db, cookie)
}

func (db DB) DeleteSession(cookie string) error {
	return deleteSession(db, cookie)
}

func (tx Tx) CreateSession(s gauth.Session) error {
	return createSession(tx, s)
}

func (tx Tx) SessionByCookie(cookie string) (gauth.Session, error) {
	return sessionByCookie(tx, cookie)
}

func (tx Tx) DeleteSession(cookie string) error {
	return deleteSession(tx, cookie)
}

func createSession(ext Ext, s gauth.Session) error {
	q := `INSERT INTO session (cookie, user_email, expires_at) VALUES ($1, $2, $3)`
	_, err := ext.Exec(q, s.Cookie, s.UserEmail, s.ExpiresAt)
	return err
}

func sessionByCookie(ext Ext, cookie string) (gauth.Session, error) {
	var session gauth.Session
	q := `SELECT * FROM session WHERE cookie = $1`
	if err := ext.Get(&session, q, cookie); err != nil {
		if err == sql.ErrNoRows {
			return gauth.Session{}, notFoundErr
		}

		return gauth.Session{}, err
	}

	return session, nil
}

func deleteSession(ext Ext, cookie string) error {
	q := `DELETE FROM session WHERE cookie = $1`
	if _, err := ext.Exec(q, cookie); err != nil {
		return err
	}

	return nil
}
