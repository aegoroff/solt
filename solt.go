package main

import (
	"github.com/spf13/afero"
	"os"
	"solt/cmd"
	"solt/cmd/out"
)

func main() {
	if err := cmd.Execute(afero.NewOsFs(), out.NewConsoleEnvironment()); err != nil {
		os.Exit(1)
	}
}
