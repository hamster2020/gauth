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
	expUser := gauth.User{Email: expEmail, PasswordHash: expHash}

	cases := []struct {
		name string

		retUserByEmailErr error

		expCheckPasswordFuncCalled bool

		expRet bool
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
			expRet:                     true,
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
				return true
			}

			ret, err := authenticate(ds, expCred, checkPasswordFunc)
			require.Equal(t, tc.expErr, err)
			require.Equal(t, tc.expRet, ret)

			require.True(t, userByEmailFuncCalled)
			require.Equal(t, tc.expCheckPasswordFuncCalled, checkPasswordFuncCalled)
		})
	}
}
