package logic

import (
	"errors"
	"testing"

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

	expToken := "token"

	cases := []struct {
		name string

		retUserByEmailErr error

		expCheckPasswordFuncCalled bool
		retCheckPasswordFunc       bool

		expNewUserTokenCalled bool
		retNewUserTokenErr    error

		expRet string
		expErr error
	}{
		{
			name:              "ds.UserByEmail error",
			retUserByEmailErr: errors.New("ds.UserByEmail error"),
			expErr:            errors.New("ds.UserByEmail error"),
		},
		{
			name:                       "ok",
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       false,
			expErr:                     errors.New("unauthorized"),
		},
		{
			name:                       "NewUserToken error",
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       true,
			expNewUserTokenCalled:      true,
			retNewUserTokenErr:         errors.New("NewUserToken error"),
			expErr:                     errors.New("NewUserToken error"),
		},
		{
			name:                       "ok",
			expCheckPasswordFuncCalled: true,
			retCheckPasswordFunc:       true,
			expNewUserTokenCalled:      true,
			expRet:                     expToken,
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

			token := mocks.NewMockToken()

			newUserTokenCalled := false
			token.NewUserTokenFunc = func(email string, roles gauth.Roles) (string, error) {
				newUserTokenCalled = true
				require.Equal(t, expEmail, email)
				require.Equal(t, expRoles, roles)
				return expToken, tc.retNewUserTokenErr
			}

			ret, err := authenticate(ds, token, expCred, checkPasswordFunc)
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expRet, ret)

			require.True(t, userByEmailFuncCalled)
			require.Equal(t, tc.expCheckPasswordFuncCalled, checkPasswordFuncCalled)
			require.Equal(t, tc.expNewUserTokenCalled, newUserTokenCalled)
		})
	}
}
