package token

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
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

func validateStandardClaims(claims jwt.StandardClaims) error {
	if claims.Subject == "" {
		return errors.New("missing subject")
	}

	now := time.Now().UTC().Unix()
	if claims.ExpiresAt < now {
		return errors.New("token expired")
	}

	if claims.IssuedAt > now {
		return errors.New("received token that was issued before current time")
	}

	if claims.NotBefore > now {
		return errors.New("received token where the not_before is before the current time")
	}

	return nil
}
