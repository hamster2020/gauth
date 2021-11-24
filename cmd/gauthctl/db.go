package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/chrismrivera/cmd"
	"github.com/hamster2020/gauth"
	"github.com/hamster2020/gauth/postgres"
)

func init() {
	app.AddCommand(createAdminCommand)
	app.AddCommand(migrateDBCommand)
}

const defaultDatabaseURL = "postgres:///gauth?sslmode=disable"

func getDatabaseUrl() string {
	url := strings.TrimSpace(os.Getenv("GAUTH_DB_URL"))
	if url == "" {
		url = defaultDatabaseURL
		fmt.Println("GAUTH_DB_URL not set, using default")
	}

	return url
}

var createAdminCommand = cmd.NewCommand("create-admin", "db", "Create a new admin user",
	func(cmd *cmd.Command) {
		cmd.AppendArg("email", "Email for the new admin user")
		cmd.Flags.String("password", "", "Password for the user")
	},

	func(cmd *cmd.Command) error {
		email := cmd.Arg("email").String()
		dbUrl := getDatabaseUrl()

		db, err := postgres.NewDB("gauth", dbUrl)
		if err != nil {
			return err
		}

		user := gauth.User{
			Email: email,
			Roles: gauth.RolesAdmin,
		}

		if _, err := db.UserByEmail(email); err == nil {
			return errors.New("user already exists with email")
		}

		password := cmd.Flag("password").String()
		if password == "" {
			fmt.Printf("Enter password for %s:", email)
			if _, err := fmt.Scan(&password); err != nil {
				return err
			}
		}

		hash, err := gauth.HashPassword(password)
		if err != nil {
			return err
		}

		user.PasswordHash = hash
		if err := db.CreateUser(user); err != nil {
			return err
		}

		fmt.Println("Admin user created")
		return nil
	},
)

var migrateDBCommand = cmd.NewCommand("migrate-db", "db", "Apply database migrations",
	func(cmd *cmd.Command) {},

	func(cmd *cmd.Command) error {
		dbUrl := getDatabaseUrl()
		db, err := postgres.NewDB("gauth", dbUrl)
		if err != nil {
			return err
		}

		version, err := db.Migrate()
		if err != nil {
			return err
		}

		fmt.Println("Migrated db to version ", version)
		return nil
	},
)
