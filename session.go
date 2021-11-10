package gauth

import (
	"time"
)

type Session struct {
	Cookie    string
	UserEmail string
	ExpiresAt time.Time
}

type SessionStore interface {
	//CRUD
	CreateSession(session Session) error
	SessionByCookie(cookie string) (Session, error)
	DeleteSession(email string) error
}
