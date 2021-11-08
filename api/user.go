package api

import (
	"errors"
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

		caller := getCallerFromRequest(req)
		if r.Roles.HasRole(gauth.RolesAdmin) && !caller.HasRole(gauth.RolesAdmin) {
			writeJSONError(crw, errors.New("only admin can create a new admin user"), http.StatusBadRequest)
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
		caller := getCallerFromRequest(req)
		if caller.Info() != email && !caller.HasRole(gauth.RolesAdmin) {
			writeJSONError(crw, errors.New("non-admin users can only look up their own user account"), http.StatusBadRequest)
			return
		}

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
		caller := getCallerFromRequest(req)
		if caller.Info() != email && !caller.HasRole(gauth.RolesAdmin) {
			writeJSONError(crw, errors.New("non-admin users can only update their own user account"), http.StatusBadRequest)
			return
		}

		if r.Roles.HasRole(gauth.RolesAdmin) && !caller.HasRole(gauth.RolesAdmin) {
			writeJSONError(crw, errors.New("non-admin users update their roles to include admin users"), http.StatusBadRequest)
			return
		}

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
		caller := getCallerFromRequest(req)
		if caller.Info() != email && !caller.HasRole(gauth.RolesAdmin) {
			writeJSONError(crw, errors.New("non-admin users can only delete their own user account"), http.StatusBadRequest)
			return
		}

		if err := logic.DeleteUser(email); err != nil {
			writeJSONError(crw, err, http.StatusInternalServerError)
			return
		}

		writeJSON(crw, nil)
	}
}
