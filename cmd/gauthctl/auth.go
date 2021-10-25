package main

import (
	"fmt"
	"net/http"
	"syscall"

	"github.com/chrismrivera/cmd"
	"github.com/hamster2020/gauth"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	app.AddCommand(authenticateCommand)
}

var authenticateCommand = cmd.NewCommand("auth", "Authenticate", "Authenticate",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email")
	},

	func(cmd *cmd.Command) error {
		fmt.Println("please provide a password")
		passwordBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			return err
		}

		cred := gauth.Credentials{
			Email:    cmd.Arg("email").String(),
			Password: string(passwordBytes),
		}

		req, err := app.makeRequest(http.MethodPost, "authenticate", cred)
		if err != nil {
			return err
		}

		if err := app.do(req, nil); err != nil {
			return err
		}

		fmt.Println("authenticated")
		return nil
	},
)
