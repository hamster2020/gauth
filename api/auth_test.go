package api

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	th := newTestHandler(t)
	expEmail := "email@test.com"
	expPassword := "P@ssword"
	expCred := gauth.Credentials{Email: expEmail, Password: expPassword}

	cases := []struct {
		name string
		body interface{}

		expCalled bool
		retPass   bool
		retErr    error

		expStatus int
		expBody   string
	}{
		{
			name:      "invalid body",
			body:      "bogus data",
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"JSON decoding failed: json: cannot unmarshal string into Go value of type gauth.Credentials"}`,
		},
		{
			name:      "logic.Autheticate error",
			body:      expCred,
			expCalled: true,
			retPass:   true,
			retErr:    errors.New("logic.Autheticate error"),
			expStatus: http.StatusInternalServerError,
			expBody:   `{"error":"logic.Autheticate error"}`,
		},
		{
			name:      "invalid credentials",
			body:      expCred,
			expCalled: true,
			retPass:   false,
			expStatus: http.StatusUnauthorized,
			expBody:   `{"error":"invalid credentials"}`,
		},
		{
			name:      "ok",
			body:      expCred,
			expCalled: true,
			retPass:   true,
			expStatus: http.StatusOK,
			expBody:   `{}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.AuthenticateFunc = func(c gauth.Credentials) (bool, error) {
				called = true
				require.Equal(t, tc.body, c)
				return tc.retPass, tc.retErr
			}

			byt := mustMarshal(t, tc.body)
			bodyReader := bytes.NewBuffer(byt)
			req := th.makeRequest(t, http.MethodPost, "/authenticate", bodyReader)

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
