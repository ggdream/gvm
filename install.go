package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v3"
)

func install(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit("gvm install <version>", 1)
	}

	var suffix string
	if runtime.GOOS == "windows" {
		suffix = "zip"
	} else {
		suffix = "tar.gz"
	}
	name := fmt.Sprintf("go%s.%s-%s.%s", cmd.Args().First(), runtime.GOOS, runtime.GOARCH, suffix)
	tempPath := filepath.Join(os.TempDir(), name)
	_, err := os.Stat(tempPath)
	if err != nil {
		if os.IsNotExist(err) {
			url := gvmhost + "/dl/" + name
			res, err := http.Get(url)
			if err != nil {
				return err
			}
			defer res.Body.Close()
			file, err := os.CreateTemp(os.TempDir(), tempPath)
			if err != nil {
				return err
			}

			bar := progressbar.DefaultBytes(
				res.ContentLength,
				fmt.Sprintf("Downloading Go%s\n", cmd.Args().First()),
			)
			if _, err = io.Copy(io.MultiWriter(file, bar), res.Body); err != nil {
				_ = file.Close()
				return err
			}

			if err = file.Close(); err != nil {
				return err
			}
			if err = os.Rename(file.Name(), tempPath); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		fmt.Printf("File %s already exists\n", tempPath)
	}

	archiver, err := NewArchiver(suffix)
	if err != nil {
		return err
	}

	err = archiver.Extract(goversions, tempPath, runtime.NumCPU()-1)
	if err != nil {
		return nil
	}

	err = os.Rename(filepath.Join(goversions, "go"), filepath.Join(goversions, cmd.Args().First()))
	if err != nil {
		return err
	}

	fmt.Printf("Installed Go%s\n", cmd.Args().First())

	return nil
}
