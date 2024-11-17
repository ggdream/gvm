package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	gvmhome    = os.Getenv("GVM_HOME")
	gvmhost    = os.Getenv("GVM_HOST")
	goroot     = runtime.GOROOT()
	goversions string
)

func init() {
	if gvmhome == "" {
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		gvmhome = filepath.Join(dir, "gvm")
	}

	if gvmhost == "" {
		// gvmhost = "https://go.dev"
		gvmhost = "https://golang.google.cn"
	}

	goversions = filepath.Join(gvmhome, "versions")
}
