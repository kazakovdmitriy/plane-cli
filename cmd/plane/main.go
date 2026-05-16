package main

import (
	"os"

	"github.com/makeplane/plane-cli/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(cli.ExitCode())
}
