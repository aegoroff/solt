package main

import (
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"io"
	"os"
	"solt/cmd"
	"solt/cmd/api"
)

func main() {
	if err := cmd.Execute(afero.NewOsFs(), newConsoleEnvironment()); err != nil {
		os.Exit(1)
	}
}

type consoleEnvironment struct{}

func newConsoleEnvironment() api.PrintEnvironment {
	return &consoleEnvironment{}
}

func (e *consoleEnvironment) NewPrinter() api.Printer {
	return api.NewPrinter(e)
}

func (*consoleEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	color.Fprintf(w, format, a...)
}

func (*consoleEnvironment) Writer() io.WriteCloser {
	return os.Stdout
}
