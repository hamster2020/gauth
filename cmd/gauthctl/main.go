package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chrismrivera/cmd"
	"github.com/hamster2020/gauth/gauthclient"
)

func getGauthURL() string {
	url := strings.TrimSpace(os.Getenv("GAUTH_WEB_URL"))
	if url == "" {
		url = "http://localhost:3000"
		fmt.Println("GAUTH_WEB_URL not set, using default")
	}

	return url
}

func main() {
	app.Description = "A command-line interface for doorman"

	if len(os.Args) >= 2 && os.Args[1] != "--help" {
		if _, ok := app.Commands[os.Args[1]]; ok {
			app.baseURL = getGauthURL()

			sessionInfo, err := readSessionInfo()
			if err != nil {
				log.Fatal(err)
			}

			app.gauthClient, err = gauthclient.NewGauthClient(
				getGauthURL(),
				sessionInfo.Email,
				"",
				sessionInfo.Token,
				sessionInfo.Cookie,
				saveSessionInfo,
			)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := app.Run(os.Args); err != nil {
		if ue, ok := err.(*cmd.UsageErr); ok {
			ue.ShowUsage()
		} else {
			fmt.Fprintln(os.Stderr, "Error:", err)
		}

		os.Exit(1)
	}
}
