package logic

import (
	"errors"
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/mocks"
	"github.com/stretchr/testify/require"
)

func TestCreatUser(t *testing.T) {
	expEmail := "test@gmail.com"
	expPassword := "P@ssword"
	expHash := "hash"
	expRoles := gauth.RolesAdmin

	expReq := gauth.UserRequest{
		Email:    expEmail,
		Password: expPassword,
		Roles:    expRoles,
	}

	expUser := gauth.User{
		Email:        expEmail,
		PasswordHash: expHash,
		Roles:        expRoles,
	}

	cases := []struct {
		name string

		retValidateEmailErr error

		expValidatePasswordCalled bool
		retValidatePasswordErr    error

		expHashPasswordFuncCalled bool
		retHashPasswordFuncErr    error

		expCreateUserCalled bool
		retCreateUserErr    error

		expErr error
	}{
		{
			name:                "emailValidator.Validate error",
			retValidateEmailErr: errors.New("emailValidator.Validate error"),
			expErr:              errors.New("emailValidator.Validate error"),
		},
		{
			name:                      "passwordValidator.Validate error",
			expValidatePasswordCalled: true,
			retValidatePasswordErr:    errors.New("passwordValidator.Validate error"),
			expErr:                    errors.New("passwordValidator.Validate error"),
		},
		{
			name:                      "hashPasswordFunc error",
			expValidatePasswordCalled: true,
			expHashPasswordFuncCalled: true,
			retHashPasswordFuncErr:    errors.New("hashPasswordFunc error"),
			expErr:                    errors.New("hashPasswordFunc error"),
		},
		{
			name:                      "ds.CreateUser error",
			expValidatePasswordCalled: true,
			expHashPasswordFuncCalled: true,
			expCreateUserCalled:       true,
			retCreateUserErr:          errors.New("ds.CreateUser error"),
			expErr:                    errors.New("ds.CreateUser error"),
		},
		{
			name:                      "ok",
			expValidatePasswordCalled: true,
			expHashPasswordFuncCalled: true,
			expCreateUserCalled:       true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			emailValidator := mocks.NewMockValidator()

			validateEmailCalled := false
			emailValidator.ValidateFunc = func(email string) error {
				validateEmailCalled = true
				require.Equal(t, expEmail, email)
				return tc.retValidateEmailErr
			}

			passwordValidator := mocks.NewMockValidator()

			validatePasswordCalled := false
			passwordValidator.ValidateFunc = func(password string) error {
				validatePasswordCalled = true
				require.Equal(t, expPassword, password)
				return tc.retValidatePasswordErr
			}

			hashPasswordFuncCalled := false
			hashPasswordFunc := func(password string) (string, error) {
				hashPasswordFuncCalled = true
				require.Equal(t, expPassword, password)
				return expHash, tc.retHashPasswordFuncErr
			}

			ds := mocks.NewMockDatastore()

			createUserCalled := false
			ds.CreateUserFunc = func(user gauth.User) error {
				createUserCalled = true
				require.Equal(t, expUser, user)
				return tc.retCreateUserErr
			}

			err := createUser(ds, emailValidator, passwordValidator, expReq, hashPasswordFunc)
			require.Equal(t, tc.expErr, err)

			require.True(t, validateEmailCalled)
			require.Equal(t, tc.expValidatePasswordCalled, validatePasswordCalled)
			require.Equal(t, tc.expHashPasswordFuncCalled, hashPasswordFuncCalled)
			require.Equal(t, tc.expCreateUserCalled, createUserCalled)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	expEmail := "test@gmail.com"
	expNewEmail := "new@gmail.com"
	expPassword := "P@ssword"
	expHash := "hash"
	expRoles := gauth.RolesAdmin

	expReq := gauth.UserRequest{
		Email:    expNewEmail,
		Password: expPassword,
		Roles:    expRoles,
	}

	expUser := gauth.User{
		Email:        expNewEmail,
		PasswordHash: expHash,
		Roles:        expRoles,
	}

	cases := []struct {
		name string

		retUserByEmailErr error

		expValidateEmailCalled bool
		retValidateEmailErr    error

		expValidatePasswordCalled bool
		retValidatePasswordErr    error

		expHashPasswordFuncCalled bool
		retHashPasswordFuncErr    error

		expUpdateUserCalled bool
		retUpdateUserErr    error

		expErr error
	}{
		{
			name:              "ds.UserByEmail error",
			retUserByEmailErr: errors.New("ds.UserByEmail error"),
			expErr:            errors.New("ds.UserByEmail error"),
		},
		{
			name:                   "emailValidator.Validate error",
			expValidateEmailCalled: true,
			retValidateEmailErr:    errors.New("emailValidator.Validate error"),
			expErr:                 errors.New("emailValidator.Validate error"),
		},
		{
			name:                      "passwordValidator.Validate error",
			expValidateEmailCalled:    true,
			expValidatePasswordCalled: true,
			retValidatePasswordErr:    errors.New("passwordValidator.Validate error"),
			expErr:                    errors.New("passwordValidator.Validate error"),
		},
		{
			name:                      "hashPasswordFunc error",
			expValidateEmailCalled:    true,
			expValidatePasswordCalled: true,
			expHashPasswordFuncCalled: true,
			retHashPasswordFuncErr:    errors.New("hashPasswordFunc error"),
			expErr:                    errors.New("hashPasswordFunc error"),
		},
		{
			name:                      "ds.UpdateUser error",
			expValidateEmailCalled:    true,
			expValidatePasswordCalled: true,
			expHashPasswordFuncCalled: true,
			expUpdateUserCalled:       true,
			retUpdateUserErr:          errors.New("ds.UpdateUser error"),
			expErr:                    errors.New("ds.UpdateUser error"),
		},
		{
			name:                      "ok",
			expValidateEmailCalled:    true,
			expValidatePasswordCalled: true,
			expHashPasswordFuncCalled: true,
			expUpdateUserCalled:       true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ds := mocks.NewMockDatastore()

			userByEmailFuncCalled := false
			ds.UserByEmailFunc = func(email string) (gauth.User, error) {
				userByEmailFuncCalled = true
				require.Equal(t, expEmail, email)
				return gauth.User{}, tc.retUserByEmailErr
			}

			emailValidator := mocks.NewMockValidator()

			validateEmailCalled := false
			emailValidator.ValidateFunc = func(email string) error {
				validateEmailCalled = true
				require.Equal(t, expNewEmail, email)
				return tc.retValidateEmailErr
			}

			passwordValidator := mocks.NewMockValidator()

			validatePasswordCalled := false
			passwordValidator.ValidateFunc = func(password string) error {
				validatePasswordCalled = true
				require.Equal(t, expPassword, password)
				return tc.retValidatePasswordErr
			}

			hashPasswordFuncCalled := false
			hashPasswordFunc := func(password string) (string, error) {
				hashPasswordFuncCalled = true
				require.Equal(t, expPassword, password)
				return expHash, tc.retHashPasswordFuncErr
			}

			updateUserCalled := false
			ds.UpdateUserFunc = func(email string, user gauth.User) error {
				updateUserCalled = true
				require.Equal(t, expEmail, email)
				require.Equal(t, expUser, user)
				return tc.retUpdateUserErr
			}

			err := updateUser(ds, emailValidator, passwordValidator, expEmail, expReq, hashPasswordFunc)
			require.Equal(t, tc.expErr, err)

			require.True(t, userByEmailFuncCalled)
			require.Equal(t, tc.expValidateEmailCalled, validateEmailCalled)
			require.Equal(t, tc.expValidatePasswordCalled, validatePasswordCalled)
			require.Equal(t, tc.expHashPasswordFuncCalled, hashPasswordFuncCalled)
			require.Equal(t, tc.expUpdateUserCalled, updateUserCalled)
		})
	}
}
