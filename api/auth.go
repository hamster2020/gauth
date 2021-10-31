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

		token, err := logic.Authenticate(r)
		if err != nil {
			writeJSONError(crw, err, http.StatusUnauthorized)
			return
		}

		writeJSON(crw, token)
	}
}
