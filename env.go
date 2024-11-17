package main

import (
	"context"
	"fmt"
	"runtime"

	"github.com/urfave/cli/v3"
)

func env(ctx context.Context, cmd *cli.Command) error {
	fmt.Printf("GOOS=%s\n", runtime.GOOS)
	fmt.Printf("GOARCH=%s\n", runtime.GOARCH)
	fmt.Printf("GOROOT=%s\n", goroot)
	fmt.Printf("GVM_HOME=%s\n", gvmhome)
	fmt.Printf("GVM_HOST=%s\n", gvmhost)

	return nil
}
