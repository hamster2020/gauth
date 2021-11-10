package logic

import (
	"errors"
	"testing"
	"time"

	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/mocks"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	expEmail := "test@gmail.com"
	expPassword := "P@ssword"
	expCred := gauth.Credentials{Email: expEmail, Password: expPassword}

	expHash := "hash"
	expRoles := gauth.RolesAdmin
	expUser := gauth.User{Email: expEmail, PasswordHash: expHash, Roles: expRoles}

	expCookie := "cookie!"

	expToken := "token"

	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	expSession := gauth.Session{Cookie: expCookie, UserEmail: expEmail, ExpiresAt: tomorrow}

	cases := []struct {
		name   string
		c      gauth.Credentials
		cookie string

		expUserByEmailCalled bool
		retUserByEmailErr    error

		expCheckPasswordFuncCalled bool
		retCheckPasswordFunc       bool

		expNewSessionFuncCalled bool
		retNewSessionFuncErr    error

		expCreateSessionCalled bool
		retCreateSessionErr    error

		expSessionByCookieCalled bool
		retSessionByCookieErr    error

		expNewUserTokenCalled bool
		retNewUserTokenErr    error

		expToken   string
		expSession gauth.Session
		expErr     error
	}{
		{
			name:                 "credentials - ds.UserByEmail error",
			c:                    expCred,
			expUserByEmailCalled: true,
			retUserByEmailErr:    errors.New("ds.UserByEmail error"),
			expErr:               errors.New("ds.UserByEmail error"),
		},
		{
			name:                       "credentials - checkPasswordFunc error",
			c:                          expCred,
			expUserByEmailCalled:       true,
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       false,
			expErr:                     errors.New("unauthorized"),
		},
		{
			name:                       "credentials - newSessionFunc error",
			c:                          expCred,
			expUserByEmailCalled:       true,
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       true,
			expNewSessionFuncCalled:    true,
			retNewSessionFuncErr:       errors.New("ds.UserByEmail error"),
			expErr:                     errors.New("ds.UserByEmail error"),
		},
		{
			name:                       "credentials - ds.CreateSession error",
			c:                          expCred,
			expUserByEmailCalled:       true,
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       true,
			expNewSessionFuncCalled:    true,
			expCreateSessionCalled:     true,
			retCreateSessionErr:        errors.New("ds.CreateSession error"),
			expErr:                     errors.New("ds.CreateSession error"),
		},
		{
			name:                     "cookie - ds.SessionByCookie error",
			cookie:                   expCookie,
			expSessionByCookieCalled: true,
			retSessionByCookieErr:    errors.New("ds.SessionByCookie error"),
			expErr:                   errors.New("ds.SessionByCookie error"),
		},
		{
			name:                     "cookie - ds.UserByEmail error",
			cookie:                   expCookie,
			expSessionByCookieCalled: true,
			expUserByEmailCalled:     true,
			retUserByEmailErr:        errors.New("ds.UserByEmail error"),
			expErr:                   errors.New("ds.UserByEmail error"),
		},
		{
			name:   "no auth provided",
			expErr: errors.New("must provide either email and password or session cookie"),
		},
		{
			name:                       "token.NewUserToken error",
			c:                          expCred,
			expUserByEmailCalled:       true,
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       true,
			expNewSessionFuncCalled:    true,
			expCreateSessionCalled:     true,
			expNewUserTokenCalled:      true,
			retNewUserTokenErr:         errors.New("token.NewUserToken error"),
			expErr:                     errors.New("token.NewUserToken error"),
		},
		{
			name:                       "credentials - ok",
			c:                          expCred,
			expUserByEmailCalled:       true,
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       true,
			expNewSessionFuncCalled:    true,
			expCreateSessionCalled:     true,
			expNewUserTokenCalled:      true,
			expToken:                   expToken,
			expSession:                 expSession,
		},
		{
			name:                     "cookie - ok",
			cookie:                   expCookie,
			expSessionByCookieCalled: true,
			expUserByEmailCalled:     true,
			expNewUserTokenCalled:    true,
			expToken:                 expToken,
			expSession:               expSession,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ds := mocks.NewMockDatastore()

			userByEmailFuncCalled := false
			ds.UserByEmailFunc = func(email string) (gauth.User, error) {
				userByEmailFuncCalled = true
				require.Equal(t, expEmail, email)
				return expUser, tc.retUserByEmailErr
			}

			checkPasswordFuncCalled := false
			checkPasswordFunc := func(password, hash string) bool {
				checkPasswordFuncCalled = true
				require.Equal(t, expPassword, password)
				require.Equal(t, expHash, hash)
				return tc.retCheckPasswordFunc
			}

			newSessionFuncCalled := false
			newSessionFunc := func(email string) (gauth.Session, error) {
				newSessionFuncCalled = true
				require.Equal(t, expEmail, email)
				return expSession, tc.retNewSessionFuncErr
			}

			createSessionCalled := false
			ds.CreateSessionFunc = func(session gauth.Session) error {
				createSessionCalled = true
				require.Equal(t, expSession, session)
				return tc.retCreateSessionErr
			}

			sessionByCookieCalled := false
			ds.SessionByCookieFunc = func(cookie string) (gauth.Session, error) {
				sessionByCookieCalled = true
				require.Equal(t, expCookie, cookie)
				return expSession, tc.retSessionByCookieErr
			}

			tokens := mocks.NewMockToken()

			newUserTokenCalled := false
			tokens.NewUserTokenFunc = func(email string, roles gauth.Roles) (string, error) {
				newUserTokenCalled = true
				require.Equal(t, expEmail, email)
				require.Equal(t, expRoles, roles)
				return expToken, tc.retNewUserTokenErr
			}

			token, session, err := authenticate(
				ds,
				tokens,
				tc.c,
				tc.cookie,
				newSessionFunc,
				checkPasswordFunc,
			)
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expToken, token)
			require.Equal(t, tc.expSession, session)

			require.Equal(t, tc.expUserByEmailCalled, userByEmailFuncCalled)
			require.Equal(t, tc.expCheckPasswordFuncCalled, checkPasswordFuncCalled)
			require.Equal(t, tc.expNewSessionFuncCalled, newSessionFuncCalled)
			require.Equal(t, tc.expCreateSessionCalled, createSessionCalled)
			require.Equal(t, tc.expSessionByCookieCalled, sessionByCookieCalled)
			require.Equal(t, tc.expNewUserTokenCalled, newUserTokenCalled)
		})
	}
}
