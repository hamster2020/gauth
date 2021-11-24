package gauth

import (
	"time"
)

type Session struct {
	Cookie    string    `db:"cookie"`
	UserEmail string    `db:"user_email"`
	ExpiresAt time.Time `db:"expires_at"`
}

type SessionStore interface {
	//CRUD
	CreateSession(session Session) error
	SessionByCookie(cookie string) (Session, error)
	DeleteSession(cookie string) error
}
