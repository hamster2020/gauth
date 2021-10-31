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

	roles := gauth.RolesAdmin
	newRoles := gauth.RolesBase

	user := gauth.User{Email: email, PasswordHash: hash, Roles: roles}
	m := NewMapStore()

	// look up user - not found
	testUser, err := m.UserByEmail(email)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testUser)

	// update user - not found
	err = m.UpdateUser(email, user)
	require.Equal(t, notFoundErr, err)

	// create user
	err = m.CreateUser(user)
	require.NoError(t, err)

	// create user - already exists
	err = m.CreateUser(user)
	require.Equal(t, userExistsErr, err)

	// look up user
	testUser, err = m.UserByEmail(email)
	require.NoError(t, err)
	require.Equal(t, user, testUser)

	// update roles
	user.Roles = newRoles
	err = m.UpdateUser(email, user)
	require.NoError(t, err)

	testUser, err = m.UserByEmail(email)
	require.NoError(t, err)
	require.Equal(t, user, testUser)

	// update hash
	user.PasswordHash = newHash
	err = m.UpdateUser(email, user)
	require.NoError(t, err)

	testUser, err = m.UserByEmail(email)
	require.NoError(t, err)
	require.Equal(t, user, testUser)

	// update email
	user.Email = newEmail
	err = m.UpdateUser(email, user)
	require.NoError(t, err)

	testUser, err = m.UserByEmail(email)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testUser)

	testUser, err = m.UserByEmail(newEmail)
	require.NoError(t, err)
	require.Equal(t, user, testUser)

	// list users
	users, err := m.Users()
	require.NoError(t, err)
	require.Len(t, users, 1)
	require.Equal(t, user, users[0])

	// delete user
	err = m.DeleteUser(newEmail)
	require.NoError(t, err)

	testUser, err = m.UserByEmail(newEmail)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testUser)
}
