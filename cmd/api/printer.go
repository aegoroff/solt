package api

import (
	"fmt"
	"io"
	"text/tabwriter"
)

type prn struct {
	tw *tabwriter.Writer
	w  io.Writer
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

func (r *prn) Writer() io.Writer {
	return r.w
}

func (r *prn) Twriter() *tabwriter.Writer {
	return r.tw
}

func (r *prn) Flush() {
	_ = r.tw.Flush()
}

func (r *prn) Tprint(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(r.tw, format, a...)
}

func (r *prn) Cprint(format string, a ...interface{}) {
	r.pe.PrintFunc(r.w, format, a...)
}
