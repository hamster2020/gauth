package main

import (
	"log"

	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/api"
	emailvalidator "github.com/hamster2020/gauth/email-validator"
	"github.com/hamster2020/gauth/logic"
	passwordvalidator "github.com/hamster2020/gauth/password-validator"
	"github.com/hamster2020/gauth/postgres"
	"github.com/hamster2020/gauth/token"
)

func main() {
	cfg, err := gauth.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration details: %v", err)
	}

	token, err := token.NewToken(cfg.AccessTokenExpMinutes)
	if err != nil {
		log.Fatalf("Failed to set up new token: %v", err)
	}

	db, err := postgres.NewDB("gauth", cfg.DBURL)
	if err != nil {
		log.Fatalf("Failed to set up postgres instance: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping postgres instance: %v", err)
	}

	version, err := db.Migrate()
	if err != nil {
		log.Fatalf("Failed to migrate postgres instance: %v", err)
	}
	log.Printf("postgres schema version is on version %d", version)

	emailValidator := emailvalidator.NewEmailValidator(cfg.EmailVerifierToken)
	passwordValidator := passwordvalidator.NewPasswordValidator(cfg.PwnedPasswordsURL)
	logic := logic.NewLogic(token, db, emailValidator, passwordValidator)

	apiHandler := api.NewAPIHandler(cfg, token, logic)

	apiHandler.ListenAndServe()
}
