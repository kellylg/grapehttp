package main

import (
	"grapehttp/client/cmd"
	cmdutil "grapehttp/client/cmd/util"
	"os"
)

func main() {
	cmd := cmd.NewFctlCommand(cmdutil.NewFactory(), os.Stdin, os.Stdout, os.Stderr)
	if cmd.Execute() != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
