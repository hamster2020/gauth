package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/mocks"
	"github.com/stretchr/testify/require"
)

func mustCreateConfig(t *testing.T) gauth.Config {
	cfg, err := gauth.NewConfig()
	require.NoError(t, err)
	require.NotZero(t, cfg)
	return cfg
}

type testHandler struct {
	server *httptest.Server
	cfg    gauth.Config
	logic  *mocks.MockLogic
}

func newTestHandler(t *testing.T) *testHandler {
	cfg := mustCreateConfig(t)
	logic := mocks.NewMockLogic()
	token := mocks.NewMockToken()
	return &testHandler{
		server: httptest.NewServer(NewAPIHandler(cfg, token, logic)),
		cfg:    cfg,
		logic:  logic,
	}
}

func (th testHandler) testURL(u string) string {
	return th.server.URL + u
}

func (th testHandler) makeRequest(t *testing.T, method, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, th.testURL(path), body)
	require.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")
	return req
}

func mustDo(t *testing.T, req *http.Request) *http.Response {
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotZero(t, resp)
	return resp
}

func mustMarshal(t *testing.T, v interface{}) []byte {
	byt, err := json.Marshal(v)
	require.NoError(t, err)
	return byt
}
