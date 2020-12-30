package api

import (
	"io"
	"text/tabwriter"
)

type prn struct {
	tw *tabwriter.Writer
	w  io.WriteCloser
	pe PrintEnvironment
}

// NewPrinter creates new Printer interface instance
func NewPrinter(pe PrintEnvironment) Printer {
	w := pe.Writer()
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	p := prn{
		tw: tw,
		w:  w,
		pe: pe,
	}
	return &p
}

func (r *prn) Writer() io.WriteCloser {
	return r.w
}

func (r *prn) Twriter() *tabwriter.Writer {
	return r.tw
}

func (r *prn) Cprint(format string, a ...interface{}) {
	r.pe.PrintFunc(r.w, format, a...)
}
