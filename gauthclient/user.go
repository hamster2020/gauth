package gauthclient

import (
	"fmt"
	"net/http"

	"github.com/hamster2020/gauth"
)

func (gc GauthClient) CreateUser(r gauth.UserRequest) error {
	req, err := gc.makeRequest(http.MethodPost, "/users", r)
	if err != nil {
		return err
	}

	return gc.doRequest(req, nil)
}

func (gc GauthClient) UserByEmail(email string) (gauth.User, error) {
	req, err := gc.makeRequest(http.MethodGet, fmt.Sprintf("/users/%s", email), nil)
	if err != nil {
		return gauth.User{}, err
	}

	var user gauth.User
	if err := gc.doRequest(req, &user); err != nil {
		return gauth.User{}, err
	}

	return user, nil
}

func (gc GauthClient) UpdateUser(email string, r gauth.UserRequest) error {
	u := fmt.Sprintf("/users/%s", email)
	req, err := gc.makeRequest(http.MethodPost, u, r)
	if err != nil {
		return err
	}

	return gc.doRequest(req, nil)
}

func (gc GauthClient) DeleteUser(email string) error {
	u := fmt.Sprintf("/users/%s", email)
	req, err := gc.makeRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}

	return gc.doRequest(req, nil)
}

func (gc GauthClient) ListUsers() ([]gauth.User, error) {
	req, err := gc.makeRequest(http.MethodGet, "/users", nil)
	if err != nil {
		return nil, err
	}

	var users []gauth.User
	if err := gc.doRequest(req, &users); err != nil {
		return nil, err
	}

	return users, nil
}
