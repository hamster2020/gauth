package api

import (
	"net/http"
)

type codeResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (crw *codeResponseWriter) WriteStatus(code int) {
	if crw.StatusCode == 0 {
		crw.StatusCode = code
		crw.ResponseWriter.WriteHeader(code)
	}
}

func (crw *codeResponseWriter) Write(v []byte) (int, error) {
	crw.WriteStatus(http.StatusOK)
	return crw.ResponseWriter.Write(v)
}
