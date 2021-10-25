package logic

import (
	"github.com/hamster2020/gauth"
)

type logic struct {
	ds                gauth.Datastore
	emailValidator    gauth.Validator
	passwordValidator gauth.Validator
}

func NewLogic(ds gauth.Datastore, emailvalidator, passwordvalidator gauth.Validator) logic {
	return logic{
		ds:                ds,
		emailValidator:    emailvalidator,
		passwordValidator: passwordvalidator,
	}
}
