package logic

import (
	"time"

	"github.com/hamster2020/gauth"
)

func newSession(email string) (gauth.Session, error) {
	cookie, err := gauth.RandomHex(32)
	if err != nil {
		return gauth.Session{}, err
	}

	exp := time.Now().UTC().Add(time.Hour)
	return gauth.Session{Cookie: cookie, UserEmail: email, ExpiresAt: exp}, nil
}
