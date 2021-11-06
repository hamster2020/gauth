package gauthclient

import (
	"crypto/rsa"
	"net/http"
)

func (gc GauthClient) PublicKey() (*rsa.PublicKey, error) {
	req, err := gc.makeRequest(http.MethodGet, "/publickey", nil)
	if err != nil {
		return nil, err
	}

	var pubKey *rsa.PublicKey
	if err := gc.doRequest(req, &pubKey); err != nil {
		return nil, err
	}

	return pubKey, nil
}
