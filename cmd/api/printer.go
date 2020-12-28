package api

import (
	"fmt"
	"github.com/gookit/color"
	"io"
	"text/tabwriter"
)

type prn struct {
	tw *tabwriter.Writer
	w  io.Writer
}

// NewPrinter creates new Printer interface instance
func NewPrinter(w io.Writer) Printer {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	p := prn{
		tw: tw,
		w:  w,
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
	color.Fprintf(r.w, format, a...)
}

func (*prn) SetColor(c color.Color) {
	_, _ = color.Set(c)
}

func (*prn) ResetColor() {
	_, _ = color.Reset()
}
