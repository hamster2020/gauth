package gauthclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/authenticator"
)

type GauthClient struct {
	client     *http.Client
	baseURL    string
	authClient gauth.Authenticator
}

func NewGauthClient(
	baseURL string,
	email string,
	password string,
	token string,
	cookie string,
	saveSessionCredentials func(email, token, cookie string) error,
) (GauthClient, error) {
	authClient, err := authenticator.NewAuthClient(baseURL, email, password, token, cookie, saveSessionCredentials)
	if err != nil {
		return GauthClient{}, err
	}

	return GauthClient{
		client:     &http.Client{Timeout: time.Second},
		baseURL:    baseURL,
		authClient: authClient,
	}, nil
}

func (gc GauthClient) makeRequest(method, path string, body interface{}) (*http.Request, error) {
	u := fmt.Sprintf("%s%s", gc.baseURL, path)
	var req *http.Request
	var err error

	if method != http.MethodGet {
		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(body); err == nil {
			req, err = http.NewRequest(method, u, &buf)
		}
	} else {
		req, err = http.NewRequest(method, u, nil)
	}
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (gc GauthClient) doRequest(req *http.Request, respBody interface{}) error {
	return gc.authClient.DoAuthenticatedRequest(req, respBody)
}
