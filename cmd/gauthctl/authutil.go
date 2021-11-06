package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

const gauthTokenInfoFilename = "gauth.token"

type TokenInfo struct {
	Email string
	Token string
}

func saveTokenInfo(email, token string) error {
	if email == "" {
		return clearTokenInfo()
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	dirPath := path.Join(usr.HomeDir, ".config")
	if err := os.MkdirAll(dirPath, os.FileMode(0755)); err != nil {
		return err
	}

	tokenPath := path.Join(dirPath, gauthTokenInfoFilename)
	f, err := os.OpenFile(tokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	tokenInfo := TokenInfo{
		Email: email,
		Token: token,
	}

	return enc.Encode(tokenInfo)
}

func readTokenInfo() (TokenInfo, error) {
	usr, err := user.Current()
	if err != nil {
		return TokenInfo{}, err
	}

	tokenPath := path.Join(usr.HomeDir, ".config", gauthTokenInfoFilename)
	f, err := os.Open(tokenPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return TokenInfo{}, err
		}

		return TokenInfo{}, nil
	}

	defer f.Close()

	dec := json.NewDecoder(f)
	var ti TokenInfo
	if err := dec.Decode(&ti); err != nil {
		err = clearTokenInfo()
		return TokenInfo{}, err
	}

	return ti, nil
}

func clearTokenInfo() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	tokenPath := path.Join(usr.HomeDir, ".config", gauthTokenInfoFilename)
	if err := os.Remove(tokenPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
