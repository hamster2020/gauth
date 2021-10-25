package main

import (
	"log"

	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/api"
	emailvalidator "github.com/hamster2020/gauth/email-validator"
	"github.com/hamster2020/gauth/logic"
	"github.com/hamster2020/gauth/mapstore"
	passwordvalidator "github.com/hamster2020/gauth/password-validator"
)

func main() {
	cfg, err := gauth.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration details: %v", err)
		return
	}

	db := mapstore.NewMapStore()
	emailValidator := emailvalidator.NewEmailValidator(cfg.EmailVerifierToken)
	passwordValidator := passwordvalidator.NewPasswordValidator(cfg.PwnedPasswordsURL)
	logic := logic.NewLogic(db, emailValidator, passwordValidator)
	apiHandler := api.NewAPIHandler(cfg, logic)

	apiHandler.ListenAndServe()
}
