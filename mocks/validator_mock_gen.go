// Code generated by mockit.
// DO NOT EDIT!

package mocks

type mockValidator struct {
	ValidateFunc func(s string) error
}

func NewMockValidator() *mockValidator {
	return &mockValidator{}
}

func (validator *mockValidator) Validate(s string) error {
	return validator.ValidateFunc(s)
}
