package gauth

import (
	"net/http"
)

const SessionCookieName = "session"

func GetSessionCookie(cookies []*http.Cookie) string {
	for _, cookie := range cookies {
		if cookie.Name == SessionCookieName {
			return cookie.Value
		}
	}

	return ""
}
