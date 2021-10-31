package main

import (
	"log"

	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/api"
	emailvalidator "github.com/hamster2020/gauth/email-validator"
	"github.com/hamster2020/gauth/logic"
	"github.com/hamster2020/gauth/mapstore"
	passwordvalidator "github.com/hamster2020/gauth/password-validator"
	"github.com/hamster2020/gauth/token"
)

func main() {
	cfg, err := gauth.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration details: %v", err)
		return
	}

	token, err := token.NewToken(cfg.AccessTokenExpMinutes)
	if err != nil {
		log.Fatalf("Failed to set up new token: %v", err)
		return
	}

	db := mapstore.NewMapStore()

	emailValidator := emailvalidator.NewEmailValidator(cfg.EmailVerifierToken)
	passwordValidator := passwordvalidator.NewPasswordValidator(cfg.PwnedPasswordsURL)
	logic := logic.NewLogic(token, db, emailValidator, passwordValidator)

	apiHandler := api.NewAPIHandler(cfg, token, logic)

	apiHandler.ListenAndServe()
}
