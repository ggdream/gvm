package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

func list(ctx context.Context, cmd *cli.Command) error {
	targetPath, err := os.Readlink(goroot)
	var currentVersion string
	if err == nil {
		currentVersion = filepath.Base(targetPath)
	}

	return filepath.Walk(goversions, func(path string, info os.FileInfo, err error) error {
		if goversions == path {
			return nil
		}

		if info.IsDir() {
			if info.Name() == currentVersion {
				fmt.Println(info.Name(), "(current)")
			} else {
				fmt.Println(info.Name())
			}

			return filepath.SkipDir
		}

		return nil
	})
}
