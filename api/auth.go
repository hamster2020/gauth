package api

import (
	"net/http"

	"github.com/hamster2020/gauth"
)

func authenticate(logic gauth.Logic) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)

		var r gauth.Credentials
		if err := parseJSONRequest(req, &r); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		cookie := gauth.GetSessionCookie(req.Cookies())
		token, session, err := logic.Authenticate(r, cookie)
		if err != nil {
			expireSessionCookie(crw, req)
			writeJSONError(crw, err, http.StatusUnauthorized)
			return
		}

		setSessionCookie(crw, req, session)
		writeJSON(crw, token)
	}
}

func logout(logic gauth.Logic) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)

		defer expireSessionCookie(crw, req)

		cookie := gauth.GetSessionCookie(req.Cookies())
		if err := logic.Logout(cookie); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		writeJSON(crw, nil)
	}
}
