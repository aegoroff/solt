package api

import (
	"io"
)

type prn struct {
	w  io.WriteCloser
	pe PrintEnvironment
}

// NewPrinter creates new Printer interface instance
func NewPrinter(pe PrintEnvironment) Printer {
	p := prn{
		w:  pe.Writer(),
		pe: pe,
	}
	return &p
}

func (r *prn) Writer() io.WriteCloser {
	return r.w
}

func (r *prn) Cprint(format string, a ...interface{}) {
	r.pe.PrintFunc(r.w, format, a...)
}
