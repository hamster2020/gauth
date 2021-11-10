package logic

import (
	"crypto/rand"
	"time"

	"github.com/hamster2020/gauth"
)

func newSession(email string) (gauth.Session, error) {
	byt := make([]byte, 32)
	if _, err := rand.Read(byt); err != nil {
		return gauth.Session{}, err
	}

	cookie := string(byt)
	exp := time.Now().UTC().Add(time.Hour)
	return gauth.Session{Cookie: cookie, UserEmail: email, ExpiresAt: exp}, nil
}
