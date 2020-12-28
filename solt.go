package main

import (
	"github.com/spf13/afero"
	"os"
	"solt/cmd"
	"solt/cmd/api"
)

func main() {
	if err := cmd.Execute(afero.NewOsFs(), api.NewConsoleEnvironment()); err != nil {
		os.Exit(1)
	}
}
