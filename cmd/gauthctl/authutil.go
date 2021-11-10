package main

import (
	"encoding/json"
	"os"
	"os/user"
	"path"
)

const gauthSessionInfoFilename = "gauth.token"

type SessionInfo struct {
	Email  string
	Token  string
	Cookie string
}

func saveSessionInfo(email, token, cookie string) error {
	if email == "" {
		return clearSessionInfo()
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	dirPath := path.Join(usr.HomeDir, ".config")
	if err := os.MkdirAll(dirPath, os.FileMode(0755)); err != nil {
		return err
	}

	infoPath := path.Join(dirPath, gauthSessionInfoFilename)
	f, err := os.OpenFile(infoPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return err
	}

	defer f.Close()

	enc := json.NewEncoder(f)
	sessionInfo := SessionInfo{
		Email:  email,
		Token:  token,
		Cookie: cookie,
	}

	return enc.Encode(sessionInfo)
}

func readSessionInfo() (SessionInfo, error) {
	usr, err := user.Current()
	if err != nil {
		return SessionInfo{}, err
	}

	infoPath := path.Join(usr.HomeDir, ".config", gauthSessionInfoFilename)
	f, err := os.Open(infoPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return SessionInfo{}, err
		}

		return SessionInfo{}, nil
	}

	defer f.Close()

	dec := json.NewDecoder(f)
	var info SessionInfo
	if err := dec.Decode(&info); err != nil {
		err = clearSessionInfo()
		return SessionInfo{}, err
	}

	return info, nil
}

func clearSessionInfo() error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	infoPath := path.Join(usr.HomeDir, ".config", gauthSessionInfoFilename)
	if err := os.Remove(infoPath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
