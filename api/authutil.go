package api

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/hamster2020/gauth"
)

var callerKey = struct{}{}

func setCallerToRequest(req *http.Request, caller gauth.Caller) {
	ctx := req.Context()
	newCtx := context.WithValue(ctx, callerKey, caller)
	*req = *req.WithContext(newCtx)
}

func getCallerFromRequest(req *http.Request) gauth.Caller {
	caller, ok := req.Context().Value(callerKey).(gauth.Caller)
	if !ok {
		return gauth.EmptyCaller{}
	}

	return caller
}

func getToken(req *http.Request) string {
	authHeader := req.Header.Get("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer: ")
}

func loadCaller(token gauth.Token, crw *codeResponseWriter, req *http.Request) error {
	tokenStr := getToken(req)
	if tokenStr == "" {
		return nil
	}

	user, err := token.VerifyUserToken(tokenStr)
	if err != nil {
		return nil
	}

	setCallerToRequest(req, user)
	return nil
}

func authByRole(h http.Handler, roles gauth.Roles) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		crw := resp.(*codeResponseWriter)

		caller := getCallerFromRequest(req)
		if !caller.HasAtLeastOneRole(roles) {
			writeJSONError(crw, errors.New("forbidden"), http.StatusForbidden)
			return
		}

		h.ServeHTTP(resp, req)
	}
}

var authAdmin = func(h http.Handler) http.HandlerFunc {
	return authByRole(h, gauth.RolesAdmin)
}

var auth = func(h http.Handler) http.HandlerFunc {
	return authByRole(h, gauth.AllRoles())
}
