package token

import "crypto/rsa"

func (t token) PublicKey() rsa.PublicKey {
	return t.accessTokenSigningKey.PublicKey
}
