package logic

import (
	"errors"

	"github.com/hamster2020/gauth"
	"golang.org/x/crypto/bcrypt"
)

func checkPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (l logic) Authenticate(c gauth.Credentials, cookie string) (string, gauth.Session, error) {
	return authenticate(l.ds, l.token, c, cookie, newSession, checkPassword)
}

func authenticate(
	ds gauth.Datastore,
	token gauth.Token,
	c gauth.Credentials,
	cookie string,
	newSessionFunc func(email string) (gauth.Session, error),
	checkPasswordFunc func(password, hash string) bool,
) (string, gauth.Session, error) {
	var user gauth.User
	var session gauth.Session
	var err error
	switch {
	// Auth via email and password
	case c.Email != "" && c.Password != "":
		user, err = ds.UserByEmail(c.Email)
		if err != nil {
			return "", gauth.Session{}, err
		}

		if pass := checkPasswordFunc(c.Password, user.PasswordHash); !pass {
			return "", gauth.Session{}, errors.New("unauthorized")
		}

		session, err = newSessionFunc(user.Email)
		if err != nil {
			return "", gauth.Session{}, err
		}

		if err := ds.CreateSession(session); err != nil {
			return "", gauth.Session{}, err
		}

	// Auth via session cookie
	case cookie != "":
		session, err = ds.SessionByCookie(cookie)
		if err != nil {
			return "", gauth.Session{}, err
		}

		user, err = ds.UserByEmail(session.UserEmail)
		if err != nil {
			return "", gauth.Session{}, err
		}

	// No Auth provided
	default:
		return "", gauth.Session{}, errors.New("must provide either email and password or session cookie")
	}

	tokenStr, err := token.NewUserToken(user.Email, user.Roles)
	if err != nil {
		return "", gauth.Session{}, err
	}

	return tokenStr, session, nil
}

func (l logic) Logout(cookie string) error {
	return l.ds.DeleteSession(cookie)
}
