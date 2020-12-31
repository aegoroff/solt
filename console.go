package main

import (
	"github.com/gookit/color"
	"io"
	"os"
	"solt/cmd/api"
)

type consoleEnvironment struct{}

func newConsoleEnvironment() api.PrintEnvironment {
	return &consoleEnvironment{}
}

func (e *consoleEnvironment) NewPrinter() (api.Printer, error) {
	return api.NewPrinter(e), nil
}

func (*consoleEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	color.Fprintf(w, format, a...)
}

func (*consoleEnvironment) Writer() io.WriteCloser {
	return os.Stdout
}
