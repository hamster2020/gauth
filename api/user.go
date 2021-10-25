package api

import (
	"net/http"

	"github.com/hamster2020/gauth"
)

func users(logic gauth.Logic) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)

		users, err := logic.Users()
		if err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		writeJSON(crw, users)
	}
}

func createUser(logic gauth.Logic) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)

		var r gauth.UserRequest
		if err := parseJSONRequest(req, &r); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		err := logic.CreateUser(r)
		if err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		writeJSON(crw, nil)
	}
}

func userByEmail(logic gauth.Logic) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)

		email := req.URL.Query().Get(":email")
		user, err := logic.UserByEmail(email)
		if err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		writeJSON(crw, user)
	}
}

func updateUser(logic gauth.Logic) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)

		var r gauth.UserRequest
		if err := parseJSONRequest(req, &r); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		email := req.URL.Query().Get(":email")
		if err := logic.UpdateUser(email, r); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		writeJSON(crw, nil)
	}
}

func deleteUser(logic gauth.Logic) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)

		email := req.URL.Query().Get(":email")
		if err := logic.DeleteUser(email); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		writeJSON(crw, nil)
	}
}
