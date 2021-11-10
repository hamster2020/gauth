package mapstore

import (
	"testing"
	"time"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func TestSessionCRUD(t *testing.T) {
	email := "email"
	cookie := "cookie"
	expires := time.Now().UTC().AddDate(0, 0, 1)

	session := gauth.Session{Cookie: cookie, UserEmail: email, ExpiresAt: expires}
	m := NewMapStore()

	// look up session - not found
	testSession, err := m.SessionByCookie(cookie)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testSession)

	// create session
	err = m.CreateSession(session)
	require.NoError(t, err)

	// create session - already exists
	err = m.CreateSession(session)
	require.Equal(t, sessionExistsErr, err)

	// look up session
	testSession, err = m.SessionByCookie(cookie)
	require.NoError(t, err)
	require.Equal(t, session, testSession)

	// delete session
	err = m.DeleteSession(cookie)
	require.NoError(t, err)

	testSession, err = m.SessionByCookie(cookie)
	require.Equal(t, notFoundErr, err)
	require.Zero(t, testSession)
}
