package passwordvalidator

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const minLen = 8

var breachedPasswordErr = errors.New("This password has been previously breached from hackers")

type PasswordValidator struct {
	client  *http.Client
	baseURL string
}

func NewPasswordValidator(baseURL string) PasswordValidator {
	return PasswordValidator{
		client:  &http.Client{Timeout: 5 * time.Second},
		baseURL: baseURL,
	}
}

func (pv PasswordValidator) isPasswordBreached(password string) error {
	sha := fmt.Sprintf("%x", sha1.Sum([]byte(password)))
	first5 := sha[:5]
	remaining := strings.ToUpper(sha[5:])
	u := fmt.Sprintf("%s/range/%s", pv.baseURL, first5)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}

	resp, err := pv.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PasswordValidator error: invalid status code of %d received from %s", resp.StatusCode, pv.baseURL)
	}

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	results := strings.Split(string(byt), `
`)
	for _, result := range results {
		r := strings.Split(result, ":")[0]
		if r == remaining {
			return breachedPasswordErr
		}
	}

	return nil
}

func (pv PasswordValidator) Validate(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be greater than %d characters", minLen)
	}

	if err := pv.isPasswordBreached(password); err != nil {
		return err
	}

	return nil
}
