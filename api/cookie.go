package api

import (
	"net/http"
	"time"

	"github.com/hamster2020/gauth"
)

func newSessionCookie(cookie string, expires time.Time, secure bool) *http.Cookie {
	return &http.Cookie{
		Name:     gauth.SessionCookieName,
		Value:    cookie,
		Path:     "/",
		Expires:  expires,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func setSessionCookie(resp http.ResponseWriter, req *http.Request, session gauth.Session) {
	secure := req.URL.Scheme == "https"
	cookie := newSessionCookie(session.Cookie, session.ExpiresAt, secure)
	http.SetCookie(resp, cookie)
}

func expireSessionCookie(resp http.ResponseWriter, req *http.Request) {
	secure := req.URL.Scheme == "https"
	cookie := newSessionCookie("", time.Now().UTC(), secure)
	http.SetCookie(resp, cookie)
}
