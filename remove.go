package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func remove(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit("gvm remove <version>", 1)
	}

	targetPath, err := os.Readlink(goroot)
	if err == nil {
		currentVersion := filepath.Base(targetPath)
		if currentVersion == cmd.Args().First() {
			err = os.RemoveAll(goroot)
			if err != nil {
				return err
			}

		}
	}

	err = os.RemoveAll(filepath.Join(goversions, cmd.Args().First()))
	if err != nil {
		return err
	}

	fmt.Printf("Removed Go%s\n", cmd.Args().First())

	return nil
}
