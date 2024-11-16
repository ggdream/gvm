package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
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
				Name:   "global",
				Usage:  "set global go version",
				Action: global,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
