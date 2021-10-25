package api

import (
	"github.com/bmizerany/pat"
)

func (h apiHandler) addRoutes() {
	router := pat.New()

	// Users
	router.Get("/users", users(h.logic))
	router.Post("/users", createUser(h.logic))
	router.Get("/users/:email", userByEmail(h.logic))
	router.Post("/users/:email", updateUser(h.logic))
	router.Del("/users/:email", deleteUser(h.logic))

	// Authenticate
	router.Post("/authenticate", authenticate(h.logic))

	h.mux.Handle("/", router)
}
