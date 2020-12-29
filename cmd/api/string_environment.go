package api

import (
	"fmt"
	"io"
	"regexp"
)

type stringEnvironment struct {
	w  io.Writer
	re *regexp.Regexp
}

func (e *stringEnvironment) NewPrinter() Printer {
	return NewPrinter(e)
}

// NewStringEnvironment creates mew plain string output environment
func NewStringEnvironment(w io.Writer) PrintEnvironment {
	return &stringEnvironment{
		w:  w,
		re: regexp.MustCompile(`<[a-zA-Z_]+>(.+?)</>`),
	}
}

func (e *stringEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(w, e.re.ReplaceAllString(s, "$1"))
}

func (e *stringEnvironment) Writer() io.Writer {
	return e.w
}
