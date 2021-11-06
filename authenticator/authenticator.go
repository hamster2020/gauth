package authenticator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/hamster2020/gauth"
)

type AuthClient struct {
	client                 *http.Client
	baseURL                string
	email                  string
	password               string
	token                  string
	saveSessionCredentials func(email, token string) error
}

func NewAuthClient(
	baseURL string,
	email string,
	password string,
	token string,
	saveSessionCredentials func(email, token string) error,
) (*AuthClient, error) {
	if baseURL == "" {
		return nil, errors.New("missing baseURL")
	}

	a := &AuthClient{
		client:                 &http.Client{Timeout: time.Second},
		baseURL:                baseURL,
		email:                  email,
		password:               password,
		token:                  token,
		saveSessionCredentials: saveSessionCredentials,
	}

	return a, nil
}

func (a *AuthClient) Authenticate() error {
	body := gauth.Credentials{
		Email:    a.email,
		Password: a.password,
	}

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseURL+"/authenticate", bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %v", resp.StatusCode)
	}

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response body: %v", err)
	}

	var token string
	if jsonErr := json.Unmarshal(byt, &token); jsonErr != nil {
		return fmt.Errorf("unable to parse response body to json: %s\n", jsonErr)
	}

	a.token = token
	if a.saveSessionCredentials != nil {
		a.saveSessionCredentials(a.email, a.token)
	}

	return nil
}

func (a *AuthClient) onAuthErr() error {
	a.token = ""
	if a.saveSessionCredentials != nil {
		a.saveSessionCredentials(a.email, a.token)
	}

	if err := a.Authenticate(); err != nil {
		a.email = ""
		a.token = ""
		if a.saveSessionCredentials != nil {
			a.saveSessionCredentials(a.email, a.token)
		}

		return err
	}

	if a.saveSessionCredentials != nil {
		a.saveSessionCredentials(a.email, a.token)
	}

	return nil
}

func (a *AuthClient) SetEmail(email string) {
	a.email = email
}

func (a *AuthClient) SetPassword(password string) {
	a.password = password
}

func (a *AuthClient) doRequestCanRetry(req *http.Request, respBody interface{}, canRetry bool) error {
	var bodyCopy io.ReadCloser
	if req.Body != nil {
		var err error
		if bodyCopy, err = req.GetBody(); err != nil {
			return fmt.Errorf("Failed to copy request body: %s", err)
		}
	}

	if a.token != "" {
		req.Header.Set("Authorization", "Bearer "+a.token)
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if canRetry {
		if resp.StatusCode == http.StatusUnauthorized {
			if err := a.onAuthErr(); err != nil {
				return err
			}

			req.Body = bodyCopy
			return a.doRequestCanRetry(req, respBody, false)
		}
	}

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Failed to read response body: %s", err)
		}
		return fmt.Errorf("Received non-JSON response: %s", strings.TrimSpace(string(body)))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if respBody != nil {
		if err := json.NewDecoder(bytes.NewBuffer(body)).Decode(respBody); err != nil {
			return err
		}
	}

	return nil
}

func (a *AuthClient) DoAuthenticatedRequest(req *http.Request, respBody interface{}) error {
	return a.doRequestCanRetry(req, respBody, true)
}
