package main

import (
	"github.com/chrismrivera/cmd"
)

func init() {
	app.AddCommand(getPublicKeyCommand)
}

var getPublicKeyCommand = cmd.NewCommand("public-key", "Key", "Lookup the public key",
	func(cmd *cmd.Command) {},

	func(cmd *cmd.Command) error {
		publicKey, err := app.gauthClient.PublicKey()
		if err != nil {
			return err
		}

		return jsonPrint(publicKey)
	},
)
