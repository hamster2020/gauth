package main

import (
	"crypto/rsa"
	"net/http"

	"github.com/chrismrivera/cmd"
)

func init() {
	app.AddCommand(getPublicKeyCommand)
}

var getPublicKeyCommand = cmd.NewCommand("public-key", "Key", "Lookup the public key",
	func(cmd *cmd.Command) {},

	func(cmd *cmd.Command) error {
		req, err := app.makeRequest(http.MethodGet, "publickey", nil)
		if err != nil {
			return err
		}

		var publicKey rsa.PublicKey
		if err := app.do(req, &publicKey); err != nil {
			return err
		}

		return jsonPrint(publicKey)
	},
)
