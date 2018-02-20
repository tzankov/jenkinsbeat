package main

import (
	"os"

	"github.com/tzankov/jenkinsbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
