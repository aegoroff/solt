package api

import (
	"io"
)

type prn struct {
	env PrintEnvironment
}

// NewPrinter creates new Printer interface instance
func NewPrinter(pe PrintEnvironment) Printer {
	return &prn{env: pe}
}

func (r *prn) Writer() io.WriteCloser {
	return r.env.Writer()
}

func (r *prn) Cprint(format string, a ...interface{}) {
	r.env.PrintFunc(r.Writer(), format, a...)
}
