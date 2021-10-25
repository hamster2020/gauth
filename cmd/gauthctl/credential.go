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
	// CRUD
	app.AddCommand(createUserCommand)
	app.AddCommand(userByEmailCommand)
	app.AddCommand(updateUserCommand)
	app.AddCommand(deleteUserCommand)

	// List
	app.AddCommand(listUsersCommand)
}

var createUserCommand = cmd.NewCommand("create-user", "User", "Create a new user",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email")
	},

	func(cmd *cmd.Command) error {
		fmt.Println("please provide a password")
		passwordBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			return err
		}

		emailAndPassword := gauth.Credentials{
			Email:    cmd.Arg("email").String(),
			Password: string(passwordBytes),
		}

		req, err := app.makeRequest(http.MethodPost, "users", emailAndPassword)
		if err != nil {
			return err
		}

		if err := app.do(req, nil); err != nil {
			return err
		}

		fmt.Println("user created")
		return nil
	},
)

var userByEmailCommand = cmd.NewCommand("get-user", "User", "Lookup an existing user by email",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email")
	},

	func(cmd *cmd.Command) error {
		u := fmt.Sprintf("users/%s", cmd.Arg("email").String())
		req, err := app.makeRequest(http.MethodGet, u, nil)
		if err != nil {
			return err
		}

		var c gauth.User
		if err := app.do(req, &c); err != nil {
			return err
		}

		return jsonPrint(c)
	},
)

var updateUserCommand = cmd.NewCommand("update-user", "User", "Update an existing user",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email")
		cmd.Flags.String("new-email", "", "New email")
		cmd.Flags.Bool("password", false, "Change password")
	},

	func(cmd *cmd.Command) error {
		emailAndPassword := gauth.Credentials{
			Email: cmd.Arg("new-email").String(),
		}

		newPassword, err := cmd.Flag("password").Bool()
		if err != nil {
			return err
		}

		if newPassword {
			fmt.Println("please provide a password")
			passwordBytes, err := terminal.ReadPassword(syscall.Stdin)
			if err != nil {
				return err
			}

			emailAndPassword.Password = string(passwordBytes)
		}

		u := fmt.Sprintf("users/%s", cmd.Arg("email").String())
		req, err := app.makeRequest(http.MethodPost, u, emailAndPassword)
		if err != nil {
			return err
		}

		if err := app.do(req, nil); err != nil {
			return err
		}

		fmt.Println("user updated")
		return nil
	},
)

var deleteUserCommand = cmd.NewCommand("delete-user", "User", "Delete a user by email",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email")
	},

	func(cmd *cmd.Command) error {
		u := fmt.Sprintf("users/%s", cmd.Arg("email").String())
		req, err := app.makeRequest(http.MethodDelete, u, nil)
		if err != nil {
			return err
		}

		if err := app.do(req, nil); err != nil {
			return err
		}

		fmt.Println("user deleted")
		return nil
	},
)

var listUsersCommand = cmd.NewCommand("list-users", "User", "List all existing users",
	func(cmd *cmd.Command) {},

	func(cmd *cmd.Command) error {
		req, err := app.makeRequest(http.MethodGet, "users", nil)
		if err != nil {
			return err
		}

		var cs []gauth.User
		if err := app.do(req, &cs); err != nil {
			return err
		}

		return jsonPrint(cs)
	},
)
