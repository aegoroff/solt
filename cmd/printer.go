package cmd

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

// newPrinter creates new printer interface instance
func newPrinter(w io.Writer) printer {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	p := prn{
		tw: tw,
		w:  w,
	}
	return &p
}

func (r *prn) writer() io.Writer {
	return r.w
}

func (r *prn) flush() {
	_ = r.tw.Flush()
}

func (r *prn) tprint(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(r.tw, format, a...)
}

func (r *prn) cprint(format string, a ...interface{}) {
	color.Fprintf(r.w, format, a...)
}

func (*prn) setColor(c color.Color) {
	_, _ = color.Set(c)
}

func (*prn) resetColor() {
	_, _ = color.Reset()
}
