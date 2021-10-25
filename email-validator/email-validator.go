package emailvalidator

import (
	"fmt"

	verifier "github.com/email-verifier/verifier-go"
)

type EmailValidator struct {
	accessToken string
}

func NewEmailValidator(token string) EmailValidator {
	return EmailValidator{accessToken: token}
}

func (ev EmailValidator) Validate(email string) error {
	valid := verifier.Verify(email, ev.accessToken)
	if !valid {
		return fmt.Errorf("The provided email '%s' is invalid", email)
	}

	return nil
}
