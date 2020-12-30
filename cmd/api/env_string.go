package api

import (
	"fmt"
	"io"
	"regexp"
)

type stringEnvironment struct {
	w  io.WriteCloser
	re *regexp.Regexp
}

func (e *stringEnvironment) NewPrinter() Printer {
	return NewPrinter(e)
}

// NewStringEnvironment creates new plain string output environment
func NewStringEnvironment(w io.WriteCloser) PrintEnvironment {
	return &stringEnvironment{
		w:  w,
		re: regexp.MustCompile(`<[a-zA-Z_=,;]+>(.+?)</>`),
	}
}

func (e *stringEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(w, e.re.ReplaceAllString(s, "$1"))
}

func (e *stringEnvironment) Writer() io.WriteCloser {
	return e.w
}
