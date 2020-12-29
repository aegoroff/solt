package api

import (
	"fmt"
	"github.com/gookit/color"
	"io"
	"os"
	"text/tabwriter"
)

type prn struct {
	tw *tabwriter.Writer
	w  io.Writer
	pe PrintEnvironment
}

// NewPrinter creates new Printer interface instance
func NewPrinter(pe PrintEnvironment) Printer {
	tw := new(tabwriter.Writer).Init(pe.Writer(), 0, 8, 4, ' ', 0)

	p := prn{
		tw: tw,
		w:  pe.Writer(),
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

type consoleEnvironment struct{}

// NewConsoleEnvironment creates mew console output environment
func NewConsoleEnvironment() PrintEnvironment {
	return &consoleEnvironment{}
}

func (c *consoleEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	color.Fprintf(w, format, a...)
}

func (c *consoleEnvironment) Writer() io.Writer {
	return os.Stdout
}

type stringEnvironment struct{ w io.Writer }

// NewConsoleEnvironment creates mew plain string output environment
func NewStringEnvironment(w io.Writer) PrintEnvironment {
	return &stringEnvironment{
		w: w,
	}
}

func (s *stringEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	_, _ = fmt.Fprintf(w, format, a...)
}

func (s *stringEnvironment) Writer() io.Writer {
	return s.w
}
