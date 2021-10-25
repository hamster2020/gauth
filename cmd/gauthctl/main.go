package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chrismrivera/cmd"
)

func getCredURL() string {
	url := strings.TrimSpace(os.Getenv("CRED_WEB_URL"))
	if url == "" {
		url = "http://localhost:3000"
		fmt.Println("CRED_WEB_URL not set, using default")
	}

	return url
}

func main() {
	app.Description = "A command-line interface for doorman"

	if len(os.Args) >= 2 && os.Args[1] != "--help" {
		if _, ok := app.Commands[os.Args[1]]; ok {
			app.baseURL = getCredURL()
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
