package main

import (
	"os"

	"github.com/enkoder/library/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
