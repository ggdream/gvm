package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	app := cli.Command{
		Name:  "gvm",
		Usage: "GVM is a tool for managing multiple Go versions",
		Commands: []*cli.Command{
			{
				Name:   "install",
				Usage:  "install go version",
				Action: install,
			},
			{
				Name:   "remove",
				Usage:  "remove go version",
				Action: remove,
			},
			{
				Name:   "global",
				Usage:  "set global go version",
				Action: global,
			},
			{
				Name:   "list",
				Usage:  "list all go versions",
				Action: list,
			},
			{
				Name:   "env",
				Usage:  "print gvm env",
				Action: env,
			},
		},
		Version: fmt.Sprintf("GVM %s (git commit %s) built on %s", version, commit, date),
	}
	if version == "" || commit == "" || date == "" {
		app.Version = "The developer custom version"
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
