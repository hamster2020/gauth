package main

import (
	"fmt"
	"strings"
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
		cmd.AppendArg("roles", "Comma-separated list of role names. Ex: 'admin,base'")
	},

	func(cmd *cmd.Command) error {
		fmt.Println("please provide a password")
		passwordBytes, err := terminal.ReadPassword(syscall.Stdin)
		if err != nil {
			return err
		}

		names := strings.Split(cmd.Arg("roles").String(), ",")
		roles, err := gauth.RolesFromNames(names)
		if err != nil {
			return err
		}

		r := gauth.UserRequest{
			Email:    cmd.Arg("email").String(),
			Password: string(passwordBytes),
			Roles:    roles,
		}

		if err := app.gauthClient.CreateUser(r); err != nil {
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
		email := cmd.Arg("email").String()
		user, err := app.gauthClient.UserByEmail(email)
		if err != nil {
			return err
		}

		return jsonPrint(user)
	},
)

var updateUserCommand = cmd.NewCommand("update-user", "User", "Update an existing user",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email")
		cmd.Flags.String("new-email", "", "New email")
		cmd.Flags.Bool("password", false, "Change password")
		cmd.Flags.String("role", "", "New comma-seaparated list of roles. Ex: 'admin,base'")
	},

	func(cmd *cmd.Command) error {
		r := gauth.UserRequest{
			Email: cmd.Flag("new-email").String(),
		}

		newPassword, err := cmd.Flag("password").Bool()
		if err != nil {
			return err
		}

		if newPassword {
			passwordBytes, err := terminal.ReadPassword(syscall.Stdin)
			if err != nil {
				return err
			}

			r.Password = string(passwordBytes)
		}

		namesStr := cmd.Flag("role").String()
		if namesStr != "" {
			names := strings.Split(namesStr, ",")
			roles, err := gauth.RolesFromNames(names)
			if err != nil {
				return err
			}
			r.Roles = roles
		}

		email := cmd.Arg("email").String()
		if err := app.gauthClient.UpdateUser(email, r); err != nil {
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
		email := cmd.Arg("email").String()
		if err := app.gauthClient.DeleteUser(email); err != nil {
			return err
		}

		fmt.Println("user deleted")
		return nil
	},
)

var listUsersCommand = cmd.NewCommand("list-users", "User", "List all existing users",
	func(cmd *cmd.Command) {},

	func(cmd *cmd.Command) error {
		users, err := app.gauthClient.ListUsers()
		if err != nil {
			return err
		}

		return jsonPrint(users)
	},
)
