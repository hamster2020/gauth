package api

import (
	"errors"
	"net/http"

	"github.com/bmizerany/pat"
)

func (h apiHandler) addRoutes() {
	router := pat.New()

	// Authenticate
	router.Post("/authenticate", authenticate(h.logic))
	router.Post("/logout", logout(h.logic))

	// Public Key
	router.Get("/publickey", publicKey(h.token))

	// Users
	router.Get("/users", authAdmin(users(h.logic)))
	router.Post("/users", createUser(h.logic))
	router.Get("/users/:email", auth(userByEmail(h.logic)))
	router.Post("/users/:email", auth(updateUser(h.logic)))
	router.Del("/users/:email", auth(deleteUser(h.logic)))

	router.NotFound = func() http.HandlerFunc {
		return func(resp http.ResponseWriter, req *http.Request) {
			crw := resp.(*codeResponseWriter)
			writeJSONError(crw, errors.New("not found"), http.StatusNotFound)
		}
	}()

	h.mux.Handle("/", router)
}
