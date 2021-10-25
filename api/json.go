package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hamster2020/gauth"
)

func parseJSONRequest(req *http.Request, v interface{}) error {
	if !strings.Contains(req.Header.Get("Content-Type"), "application/json") {
		return errors.New("Content-Type must be application/json")
	}

	if err := json.NewDecoder(req.Body).Decode(&v); err != nil {
		return fmt.Errorf("JSON decoding failed: %s", err)
	}

	return nil
}

func writeJSON(crw *codeResponseWriter, v interface{}) {
	crw.Header().Set("Content-Type", "application/json; charset=utf-8")

	if v == nil {
		v = struct{}{}
	}

	jsonBytes, err := json.Marshal(v)
	if err != nil {
		crw.WriteStatus(http.StatusInternalServerError)
		return
	}

	if _, err := crw.Write(jsonBytes); err != nil {
		crw.WriteStatus(http.StatusInternalServerError)
		return
	}
}

func writeJSONError(crw *codeResponseWriter, err error, statusCode int) {
	crw.WriteStatus(statusCode)
	writeJSON(crw, gauth.APIError{Err: err.Error()})
}
