package mapstore

import (
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func TestUserCRUD(t *testing.T) {
	email := "email"
	newEmail := "new email"

	hash := "hash"
	newHash := "new hash"

	c := gauth.User{Email: email, PasswordHash: hash}
	m := NewMapStore()

	// look up user - not found
	testC, err := m.UserByEmail(email)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testC)

	// update user - not found
	err = m.UpdateUser(email, c)
	require.Equal(t, notFoundErr, err)

	// create user
	err = m.CreateUser(c)
	require.NoError(t, err)

	// create user - already exists
	err = m.CreateUser(c)
	require.Equal(t, userExistsErr, err)

	// look up user
	testC, err = m.UserByEmail(email)
	require.NoError(t, err)
	require.Equal(t, c, testC)

	// update hash
	c.PasswordHash = newHash
	err = m.UpdateUser(email, c)
	require.NoError(t, err)

	testC, err = m.UserByEmail(email)
	require.NoError(t, err)
	require.Equal(t, c, testC)

	// update email
	c.Email = newEmail
	err = m.UpdateUser(email, c)
	require.NoError(t, err)

	testC, err = m.UserByEmail(email)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testC)

	testC, err = m.UserByEmail(newEmail)
	require.NoError(t, err)
	require.Equal(t, c, testC)

	// list users
	cs, err := m.Users()
	require.NoError(t, err)
	require.Len(t, cs, 1)
	require.Equal(t, c, cs[0])

	// delete user
	err = m.DeleteUser(newEmail)
	require.NoError(t, err)

	testC, err = m.UserByEmail(email)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testC)
}
