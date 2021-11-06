package main

import (
	"fmt"
	"syscall"

	"github.com/chrismrivera/cmd"
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

		email := cmd.Arg("email").String()
		password := string(passwordBytes)
		if err := app.gauthClient.Authenticate(email, password); err != nil {
			return err
		}

		fmt.Println("authenticated")
		return nil
	},
)
