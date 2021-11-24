package postgres

import (
	"database/sql"
	"errors"

	"github.com/hamster2020/gauth"
)

var notFoundErr = errors.New("not found")

func (db DB) CreateUser(u gauth.User) error {
	return createUser(db, u)
}

func (db DB) UserByEmail(email string) (gauth.User, error) {
	return userByEmail(db, email)
}

func (db DB) UpdateUser(email string, u gauth.User) error {
	return updateUser(db, email, u)
}

func (db DB) DeleteUser(email string) error {
	return deleteUser(db, email)
}

func (db DB) Users() ([]gauth.User, error) {
	return users(db)
}

func (tx Tx) CreateUser(u gauth.User) error {
	return createUser(tx, u)
}

func (tx Tx) UserByEmail(email string) (gauth.User, error) {
	return userByEmail(tx, email)
}

func (tx Tx) UpdateUser(email string, u gauth.User) error {
	return updateUser(tx, email, u)
}

func (tx Tx) DeleteUser(email string) error {
	return deleteUser(tx, email)
}

func (tx Tx) Users() ([]gauth.User, error) {
	return users(tx)
}

func createUser(ext Ext, u gauth.User) error {
	q := `INSERT INTO "user" (email, password_hash, roles) VALUES ($1, $2, $3)`
	_, err := ext.Exec(q, u.Email, u.PasswordHash, u.Roles)
	return err
}

func userByEmail(ext Ext, email string) (gauth.User, error) {
	var user gauth.User
	q := `SELECT * FROM "user" WHERE email = $1`
	if err := ext.Get(&user, q, email); err != nil {
		if err == sql.ErrNoRows {
			return gauth.User{}, notFoundErr
		}

		return gauth.User{}, err
	}

	return user, nil
}

func updateUser(ext Ext, email string, u gauth.User) error {
	_, err := userByEmail(ext, email)
	if err != nil {
		return err
	}

	if err := deleteUser(ext, email); err != nil {
		return nil
	}

	return createUser(ext, u)
}

func deleteUser(ext Ext, email string) error {
	q := `DELETE FROM "user" WHERE email = $`
	if _, err := ext.Exec(q, email); err != nil {
		return err
	}

	return nil
}

func users(ext Ext) ([]gauth.User, error) {
	var users []gauth.User
	q := `SELECT * FROM "user"`
	if err := ext.Select(&users, q); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return users, nil
}
