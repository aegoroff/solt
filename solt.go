package main

import (
	"github.com/spf13/afero"
	"os"
	"solt/cmd"
)

func main() {
	if err := cmd.Execute(afero.NewOsFs(), os.Stdout); err != nil {
		os.Exit(1)
	}
}
