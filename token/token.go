package token

import (
	"crypto/rand"
	"crypto/rsa"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type token struct {
	accessTokenExpMinutes uint8
	accessTokenSigningKey *rsa.PrivateKey
}

func NewToken(
	accessTokenExpMinutes uint8,
) (token, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return token{}, err
	}

	return token{
		accessTokenExpMinutes: accessTokenExpMinutes,
		accessTokenSigningKey: privateKey,
	}, nil
}

func newStandardClaims(subject string, exp time.Duration) jwt.StandardClaims {
	now := time.Now().UTC()
	return jwt.StandardClaims{
		Subject:   subject,
		ExpiresAt: now.Add(exp).Unix(),
		IssuedAt:  now.Unix(),
		NotBefore: now.Unix(),
	}
}
