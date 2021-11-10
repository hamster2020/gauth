package mapstore

import (
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func TestMapstoreUserCRUD(t *testing.T) {
	m := NewMapStore()

	key := "key"
	value := gauth.User{Email: "old"}
	newValue := gauth.User{Email: "new"}

	v, err := m.getUser(key)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, v)

	m.setUser(key, value)

	v, err = m.getUser(key)
	require.NoError(t, err)
	require.Equal(t, value, v)

	m.setUser(key, newValue)

	v, err = m.getUser(key)
	require.NoError(t, err)
	require.Equal(t, newValue, v)

	m.deleteUser(key)

	v, err = m.getUser(key)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, v)
}

func TestMapstoreSessionCRUD(t *testing.T) {
	m := NewMapStore()

	key := "key"
	value := gauth.Session{Cookie: "old"}
	newValue := gauth.Session{Cookie: "new"}

	v, err := m.getSession(key)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, v)

	m.setSession(key, value)

	v, err = m.getSession(key)
	require.NoError(t, err)
	require.Equal(t, value, v)

	m.setSession(key, newValue)

	v, err = m.getSession(key)
	require.NoError(t, err)
	require.Equal(t, newValue, v)

	m.deleteSession(key)

	v, err = m.getSession(key)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, v)
}
