package api

import (
	"net/http"

	"github.com/hamster2020/gauth"
)

func publicKey(token gauth.Token) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		crw := w.(*codeResponseWriter)
		writeJSON(crw, token.PublicKey())
	}
}
