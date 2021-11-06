package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chrismrivera/cmd"
	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/gauthclient"
)

type gauthCtl struct {
	*cmd.App
	*http.Client
	gauthClient gauthclient.GauthClient
	baseURL     string
}

var app = gauthCtl{
	App:    cmd.NewApp(),
	Client: &http.Client{Timeout: 30 * time.Second},
}

func (g gauthCtl) credURL(path string) string {
	return fmt.Sprintf("%s/%s", g.baseURL, path)
}

func (g gauthCtl) makeRequest(method string, path string, body interface{}) (*http.Request, error) {
	bodyBuffer := &bytes.Buffer{}
	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		bodyBuffer = bytes.NewBuffer(bodyJson)
	}

	req, err := http.NewRequest(method, app.credURL(path), bodyBuffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (g gauthCtl) do(req *http.Request, v interface{}) error {
	resp, err := app.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr gauth.APIError
		if jsonErr := parseResponseBodyJSON(resp.Body, &apiErr); jsonErr != nil {
			return jsonErr
		}

		return fmt.Errorf("Invalid response: %d - %v\n", resp.StatusCode, apiErr)
	}

	if v != nil {
		if err := parseResponseBodyJSON(resp.Body, v); err != nil {
			return err
		}
	}

	return nil
}
