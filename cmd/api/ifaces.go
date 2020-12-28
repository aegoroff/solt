package api

import (
	"io"
	"text/tabwriter"
)

// Matcher defines string matcher interface
type Matcher interface {
	// Match do string matching to several patterns
	Match(s string) bool
}

type Printer interface {
	Writer() io.Writer

	Twriter() *tabwriter.Writer

	Flush()

	// tprint prints using tab writer
	Tprint(format string, a ...interface{})

	// cprint prints data with suppport colorizing
	Cprint(format string, a ...interface{})
}

type Executor interface {
	Execute() error
}
