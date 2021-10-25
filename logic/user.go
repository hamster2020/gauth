package logic

import (
	"github.com/hamster2020/gauth"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (l logic) CreateUser(r gauth.UserRequest) error {
	return createUser(
		l.ds,
		l.emailValidator,
		l.passwordValidator,
		r,
		hashPassword,
	)
}

func createUser(
	ds gauth.Datastore,
	eValidator gauth.Validator,
	pValidator gauth.Validator,
	r gauth.UserRequest,
	hashPasswordFunc func(password string) (string, error),
) error {
	if err := eValidator.Validate(r.Email); err != nil {
		return err
	}

	if err := pValidator.Validate(r.Password); err != nil {
		return err
	}

	hash, err := hashPasswordFunc(r.Password)
	if err != nil {
		return err
	}

	user := gauth.User{
		Email:        r.Email,
		PasswordHash: string(hash),
		Roles:        r.Roles,
	}

	return ds.CreateUser(user)
}

func (l logic) UserByEmail(email string) (gauth.User, error) {
	return l.ds.UserByEmail(email)
}

func (l logic) UpdateUser(oldEmail string, r gauth.UserRequest) error {
	return updateUser(
		l.ds,
		l.emailValidator,
		l.passwordValidator,
		oldEmail,
		r,
		hashPassword,
	)
}

func updateUser(
	ds gauth.Datastore,
	eValidator gauth.Validator,
	pValidator gauth.Validator,
	email string,
	r gauth.UserRequest,
	hashPasswordFunc func(password string) (string, error),
) error {
	user, err := ds.UserByEmail(email)
	if err != nil {
		return err
	}

	if r.Email != "" {
		if err := eValidator.Validate(r.Email); err != nil {
			return err
		}

		user.Email = r.Email
	}

	if r.Password != "" {
		if err := pValidator.Validate(r.Password); err != nil {
			return err
		}

		user.PasswordHash, err = hashPasswordFunc(r.Password)
		if err != nil {
			return err
		}
	}

	if r.Roles != 0 {
		user.Roles = r.Roles
	}

	return ds.UpdateUser(email, user)
}

func (l logic) DeleteUser(email string) error {
	return l.ds.DeleteUser(email)
}

func (l logic) Users() ([]gauth.User, error) {
	return l.ds.Users()
}
