package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func global(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit("gvm global <version>", 1)
	}

	src := filepath.Join(gvmhome, "versions", cmd.Args().First())
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.RemoveAll(goroot)
	if err != nil {
		return err
	}
	err = os.Symlink(src, goroot)
	if err != nil {
		return err
	}

	fmt.Printf("Switched global version to Go%s\n", cmd.Args().First())

	return nil
}
