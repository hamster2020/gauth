package api

import (
	"errors"
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

		pass, err := logic.Authenticate(r)
		if err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		if !pass {
			writeJSONError(crw, errors.New("invalid credentials"), http.StatusUnauthorized)
			return
		}

		writeJSON(crw, nil)
	}
}
