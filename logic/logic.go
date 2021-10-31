package logic

import (
	"github.com/hamster2020/gauth"
)

type logic struct {
	token             gauth.Token
	ds                gauth.Datastore
	emailValidator    gauth.Validator
	passwordValidator gauth.Validator
}

func NewLogic(
	token gauth.Token,
	ds gauth.Datastore,
	emailvalidator gauth.Validator,
	passwordvalidator gauth.Validator,
) logic {
	return logic{
		token:             token,
		ds:                ds,
		emailValidator:    emailvalidator,
		passwordValidator: passwordvalidator,
	}
}
