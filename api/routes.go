package api

import (
	"github.com/bmizerany/pat"
)

func (h apiHandler) addRoutes() {
	router := pat.New()

	// Authenticate
	router.Post("/authenticate", authenticate(h.logic))

	// Public Key
	router.Get("/publickey", publicKey(h.token))

	// Users
	router.Get("/users", users(h.logic))
	router.Post("/users", createUser(h.logic))
	router.Get("/users/:email", userByEmail(h.logic))
	router.Post("/users/:email", updateUser(h.logic))
	router.Del("/users/:email", deleteUser(h.logic))

	h.mux.Handle("/", router)
}
