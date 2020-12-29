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

// NewStringEnvironment creates mew plain string output environment
func NewStringEnvironment(w io.Writer) PrintEnvironment {
	return &stringEnvironment{
		w:  w,
		re: regexp.MustCompile(`<[a-zA-Z]+>(.+?)</>`),
	}
}

func (e *stringEnvironment) PrintFunc(w io.Writer, format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	_, _ = fmt.Fprintf(w, e.re.ReplaceAllString(s, "$1"))
}

func (s *stringEnvironment) Writer() io.Writer {
	return s.w
}
