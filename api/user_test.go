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
	expUsers := []gauth.User{{Email: "test@email.com", PasswordHash: "hash", Roles: gauth.RolesAdmin}}

	cases := []struct {
		name string

		retErr error

		expStatus int
		expBody   string
	}{
		{
			name:      "logic.Users error",
			retErr:    errors.New("logic.Users error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.Users error"}`,
		},
		{
			name:      "ok",
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

			req := th.makeRequest(t, http.MethodGet, "/users", nil)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.True(t, called)
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
		name string
		body interface{}

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "invalid body",
			body:      "bogus data",
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"JSON decoding failed: json: cannot unmarshal string into Go value of type gauth.UserRequest"}`,
		},
		{
			name:      "logic.CreateUser error",
			body:      expReq,
			expCalled: true,
			retErr:    errors.New("logic.CreateUser error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.CreateUser error"}`,
		},
		{
			name:      "ok",
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
			req := th.makeRequest(t, http.MethodPost, "/users", BodyReader)
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
		name string

		retErr error

		expStatus int
		expBody   string
	}{
		{
			name:      "logic.UserByEmail error",
			retErr:    errors.New("logic.UserByEmail error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.UserByEmail error"}`,
		},
		{
			name:      "ok",
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

			req := th.makeRequest(t, http.MethodGet, fmt.Sprintf("/users/%s", expEmail), nil)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.True(t, called)
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
		name string
		body interface{}

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "invalid body",
			body:      "bogus data",
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"JSON decoding failed: json: cannot unmarshal string into Go value of type gauth.UserRequest"}`,
		},
		{
			name:      "logic.UpdateUser error",
			body:      expReq,
			expCalled: true,
			retErr:    errors.New("logic.UpdateUser error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.UpdateUser error"}`,
		},
		{
			name:      "ok",
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
			req := th.makeRequest(t, http.MethodPost, fmt.Sprintf("/users/%s", expEmail), BodyReader)
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
		name string

		retErr error

		expStatus int
		expBody   string
	}{
		{
			name:      "logic.DeleteUser error",
			retErr:    errors.New("logic.DeleteUser error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.DeleteUser error"}`,
		},
		{
			name:      "ok",
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

			req := th.makeRequest(t, http.MethodDelete, fmt.Sprintf("/users/%s", expEmail), nil)
			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))

			require.True(t, called)
		})
	}
}
