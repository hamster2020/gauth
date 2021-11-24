package main

import (
	"fmt"
	"syscall"

	"github.com/chrismrivera/cmd"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	app.AddCommand(authenticateCommand)
	app.AddCommand(logoutCommand)
	app.AddCommand(whoamiCommand)
}

var authenticateCommand = cmd.NewCommand("auth", "Auth", "Authenticate",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email")
		cmd.Flags.String("password", "", "Password for the user")
	},

	func(cmd *cmd.Command) error {
		password := cmd.Flag("password").String()
		if password == "" {
			fmt.Println("please provide a password")
			passwordBytes, err := terminal.ReadPassword(syscall.Stdin)
			if err != nil {
				return err
			}

			password = string(passwordBytes)
		}

		email := cmd.Arg("email").String()
		if err := app.gauthClient.Authenticate(email, password); err != nil {
			return err
		}

		fmt.Println("authenticated")
		return nil
	},
)

var logoutCommand = cmd.NewCommand("logout", "Auth", "Logout",
	func(cmd *cmd.Command) {},

	func(cmd *cmd.Command) error {
		if err := app.gauthClient.Logout(); err != nil {
			return err
		}

		fmt.Println("logged out")
		return nil
	},
)

var whoamiCommand = cmd.NewCommand("whoami", "Auth", "Display current user",
	func(cmd *cmd.Command) {},

	func(cmd *cmd.Command) error {
		sessionInfo, err := readSessionInfo()
		if err != nil {
			return err
		}

		fmt.Println("Email:", sessionInfo.Email)
		fmt.Println("Token:", sessionInfo.Token)
		fmt.Println("Cookie:", sessionInfo.Cookie)
		return nil
	},
)
