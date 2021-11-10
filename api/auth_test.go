package api

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/hamster2020/gauth"
	"github.com/stretchr/testify/require"
)

func TestAuthenticate(t *testing.T) {
	th := newTestHandler(t)
	expEmail := "email@test.com"
	expPassword := "P@ssword"
	expCred := gauth.Credentials{Email: expEmail, Password: expPassword}

	expToken := "token"

	expCookieStr := "cookie"
	tomorrow := time.Now().UTC().AddDate(0, 0, 1)
	expCookie := &http.Cookie{
		Name:     "session",
		Value:    expCookieStr,
		Path:     "/",
		Expires:  tomorrow,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	expSession := gauth.Session{Cookie: expCookieStr, UserEmail: expEmail, ExpiresAt: tomorrow}

	cases := []struct {
		name   string
		cookie *http.Cookie
		body   interface{}

		expCalled bool
		retErr    error

		expStatus int
		expBody   string
		expCookie string
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
			retErr:    errors.New("logic.Autheticate error"),
			expStatus: http.StatusUnauthorized,
			expBody:   `{"error":"logic.Autheticate error"}`,
		},
		{
			name:      "password - ok",
			body:      expCred,
			expCalled: true,
			expStatus: http.StatusOK,
			expBody:   `"token"`,
			expCookie: expCookieStr,
		},
		{
			name:      "cookie - ok",
			cookie:    expCookie,
			body:      gauth.Credentials{},
			expCalled: true,
			expStatus: http.StatusOK,
			expBody:   `"token"`,
			expCookie: expCookieStr,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.AuthenticateFunc = func(c gauth.Credentials, cookie string) (string, gauth.Session, error) {
				called = true
				require.Equal(t, tc.body, c)
				return expToken, expSession, tc.retErr
			}

			byt := mustMarshal(t, tc.body)
			bodyReader := bytes.NewBuffer(byt)
			req := th.makeRequest(t, http.MethodPost, "/authenticate", bodyReader, "")
			if tc.cookie != nil {
				req.AddCookie(tc.cookie)
			}

			resp := mustDo(t, req)
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)

			cookie := gauth.GetSessionCookie(resp.Cookies())

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expBody, string(body))
			require.Equal(t, tc.expCookie, cookie)

			require.Equal(t, tc.expCalled, called)
		})
	}
}

func TestLogout(t *testing.T) {
	th := newTestHandler(t)

	expCookieStr := "cookie"
	now := time.Now().UTC()
	tomorrow := now.AddDate(0, 0, 1)
	expCookie := newSessionCookie(expCookieStr, tomorrow, false)

	cases := []struct {
		name string

		expCalled bool
		retErr    error

		expStatus int
		expCookie string
	}{
		{
			name:      "logic.Logout error",
			expCalled: true,
			retErr:    errors.New("logic.Logout error"),
			expStatus: http.StatusInternalServerError,
			expCookie: "",
		},
		{
			name:      "ok",
			expCalled: true,
			expStatus: http.StatusOK,
			expCookie: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := false
			th.logic.LogoutFunc = func(cookie string) error {
				called = true
				require.Equal(t, expCookieStr, cookie)
				return tc.retErr
			}

			req := th.makeRequest(t, http.MethodPost, "/logout", nil, "")
			req.AddCookie(expCookie)

			resp := mustDo(t, req)
			defer resp.Body.Close()

			cookie := gauth.GetSessionCookie(resp.Cookies())

			require.Equal(t, tc.expStatus, resp.StatusCode)
			require.Equal(t, tc.expCookie, cookie)

			require.Equal(t, tc.expCalled, called)
		})
	}
}
