package gauth

type Validator interface {
	Validate(s string) error
}
