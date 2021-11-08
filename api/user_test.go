package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func TestUsers(t *testing.T) {
	th := newTestHandler(t)
	expEmail := "test@email.com"
	expHash := "hash"
	expRoles := gauth.RolesAdmin
	expUsers := []gauth.User{{Email: expEmail, PasswordHash: expHash, Roles: expRoles}}

	cases := []struct {
		name  string
		token string

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "non-admin",
			token:     th.newToken(t, expEmail, gauth.RolesBase),
			expStatus: http.StatusForbidden,
			expBody:   `{"error":"forbidden"}`,
		},
		{
			name:      "logic.Users error",
			token:     th.newToken(t, expEmail, expRoles),
			expCalled: true,
			retErr:    errors.New("logic.Users error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.Users error"}`,
		},
		{
			name:      "ok",
			token:     th.newToken(t, expEmail, expRoles),
			expCalled: true,
			expStatus: http.StatusOK,
			expBody:   `[{"email":"test@email.com","roles":1}]`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.UsersFunc = func() ([]gauth.User, error) {
				called = true
				return expUsers, tc.retErr
			}

			req := th.makeRequest(t, http.MethodGet, "/users", nil, tc.token)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.Equal(t, tc.expCalled, called)
		})
	}
}

func TestCreateUser(t *testing.T) {
	th := newTestHandler(t)
	expEmail := "email@test.com"
	expPassword := "P@ssword"
	expRoles := gauth.RolesAdmin
	expReq := gauth.UserRequest{
		Email:    expEmail,
		Password: expPassword,
		Roles:    expRoles,
	}

	cases := []struct {
		name  string
		token string
		body  interface{}

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "invalid body",
			token:     th.newToken(t, expEmail, expRoles),
			body:      "bogus data",
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"JSON decoding failed: json: cannot unmarshal string into Go value of type gauth.UserRequest"}`,
		},
		{
			name:      "non-admin create admin user",
			token:     th.newToken(t, expEmail, gauth.RolesBase),
			body:      expReq,
			expStatus: http.StatusBadRequest,
			expBody:   `{"error":"only admin can create a new admin user"}`,
		},
		{
			name:      "logic.CreateUser error",
			token:     th.newToken(t, expEmail, expRoles),
			body:      expReq,
			expCalled: true,
			retErr:    errors.New("logic.CreateUser error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.CreateUser error"}`,
		},
		{
			name:      "ok",
			token:     th.newToken(t, expEmail, expRoles),
			body:      expReq,
			expCalled: true,
			expStatus: http.StatusOK,
			expBody:   `{}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.CreateUserFunc = func(r gauth.UserRequest) error {
				called = true
				require.Equal(t, expReq, r)
				return tc.retErr
			}

			byt := mustMarshal(t, tc.body)
			BodyReader := bytes.NewBuffer(byt)
			req := th.makeRequest(t, http.MethodPost, "/users", BodyReader, tc.token)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.Equal(t, tc.expCalled, called)
		})
	}
}

func TestUserByEmail(t *testing.T) {
	th := newTestHandler(t)
	expEmail := "email@test.com"
	expHash := "hash"
	expRoles := gauth.RolesAdmin
	expUser := gauth.User{Email: expEmail, PasswordHash: expHash, Roles: expRoles}

	cases := []struct {
		name  string
		token string

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "no token",
			expStatus: http.StatusForbidden,
			expBody:   `{"error":"forbidden"}`,
		},
		{
			name:      "non-admin user lookup different user",
			token:     th.newToken(t, "diff@email.com", gauth.RolesBase),
			expStatus: http.StatusBadRequest,
			expBody:   `{"error":"non-admin users can only look up their own user account"}`,
		},
		{
			name:      "logic.UserByEmail error",
			token:     th.newToken(t, expEmail, expRoles),
			expCalled: true,
			retErr:    errors.New("logic.UserByEmail error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.UserByEmail error"}`,
		},
		{
			name:      "ok",
			token:     th.newToken(t, expEmail, expRoles),
			expCalled: true,
			expStatus: http.StatusOK,
			expBody:   `{"email":"email@test.com","roles":1}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.UserByEmailFunc = func(email string) (gauth.User, error) {
				called = true
				require.Equal(t, expEmail, email)
				return expUser, tc.retErr
			}

			req := th.makeRequest(t, http.MethodGet, fmt.Sprintf("/users/%s", expEmail), nil, tc.token)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.Equal(t, tc.expCalled, called)
		})
	}
}

func TestUpdateUser(t *testing.T) {
	th := newTestHandler(t)
	expEmail := "email@test.com"
	expNewEmail := "new_email@test.com"
	expNewPassword := "P@ssword"
	expNewRoles := gauth.RolesBase

	expReq := gauth.UserRequest{
		Email:    expNewEmail,
		Password: expNewPassword,
		Roles:    expNewRoles,
	}

	cases := []struct {
		name  string
		token string
		body  interface{}

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "no token",
			expStatus: http.StatusForbidden,
			expBody:   `{"error":"forbidden"}`,
		},
		{
			name:      "invalid body",
			token:     th.newToken(t, expEmail, gauth.RolesAdmin),
			body:      "bogus data",
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"JSON decoding failed: json: cannot unmarshal string into Go value of type gauth.UserRequest"}`,
		},
		{
			name:      "non-admin update different account",
			token:     th.newToken(t, "diff@email.com", gauth.RolesBase),
			expStatus: http.StatusBadRequest,
			expBody:   `{"error":"non-admin users can only update their own user account"}`,
		},
		{
			name:  "non-admin update to admin role",
			token: th.newToken(t, expEmail, gauth.RolesBase),
			body: func() gauth.UserRequest {
				ret := expReq
				ret.Roles = gauth.RolesAdmin
				return ret
			}(),
			expStatus: http.StatusBadRequest,
			expBody:   `{"error":"non-admin users update their roles to include admin users"}`,
		},
		{
			name:      "logic.UpdateUser error",
			token:     th.newToken(t, expEmail, gauth.RolesAdmin),
			body:      expReq,
			expCalled: true,
			retErr:    errors.New("logic.UpdateUser error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.UpdateUser error"}`,
		},
		{
			name:      "ok",
			token:     th.newToken(t, expEmail, gauth.RolesAdmin),
			body:      expReq,
			expCalled: true,
			expStatus: http.StatusOK,
			expBody:   `{}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.UpdateUserFunc = func(email string, r gauth.UserRequest) error {
				called = true
				require.Equal(t, expEmail, email)
				require.Equal(t, expReq, r)
				return tc.retErr
			}

			byt := mustMarshal(t, tc.body)
			BodyReader := bytes.NewBuffer(byt)
			req := th.makeRequest(t, http.MethodPost, fmt.Sprintf("/users/%s", expEmail), BodyReader, tc.token)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.Equal(t, tc.expCalled, called)
		})
	}
}

func TestDeleteUser(t *testing.T) {
	th := newTestHandler(t)
	expEmail := "email@test.com"

	cases := []struct {
		name  string
		token string

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "no token",
			expStatus: http.StatusForbidden,
			expBody:   `{"error":"forbidden"}`,
		},
		{
			name:      "non-admin delete different user account",
			token:     th.newToken(t, "diff@email.com", gauth.RolesBase),
			expStatus: http.StatusBadRequest,
			expBody:   `{"error":"non-admin users can only delete their own user account"}`,
		},
		{
			name:      "logic.DeleteUser error",
			token:     th.newToken(t, expEmail, gauth.RolesBase),
			expCalled: true,
			retErr:    errors.New("logic.DeleteUser error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.DeleteUser error"}`,
		},
		{
			name:      "ok",
			token:     th.newToken(t, expEmail, gauth.RolesBase),
			expCalled: true,
			expStatus: http.StatusOK,
			expBody:   `{}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.DeleteUserFunc = func(email string) error {
				called = true
				require.Equal(t, expEmail, email)
				return tc.retErr
			}

			req := th.makeRequest(t, http.MethodDelete, fmt.Sprintf("/users/%s", expEmail), nil, tc.token)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.Equal(t, tc.expCalled, called)
		})
	}
}
