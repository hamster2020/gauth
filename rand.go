package gauth

import (
	"crypto/rand"
	"encoding/hex"
)

func RandomHex(size int) (string, error) {
	byt := make([]byte, size)
	if _, err := rand.Read(byt); err != nil {
		return "", err
	}

	return hex.EncodeToString(byt), nil
}
