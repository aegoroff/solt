package main

import (
	"github.com/spf13/afero"
	"os"
	"solt/cmd"
	"solt/cmd/fw"
)

func main() {
	if err := cmd.Execute(afero.NewOsFs(), fw.NewConsoleEnvironment()); err != nil {
		os.Exit(1)
	}
}
