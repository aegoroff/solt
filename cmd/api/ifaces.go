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

// PrintEnvironment represents concrete printing environment abstraction
type PrintEnvironment interface {
	// PrintFunc represents printing function implementation
	PrintFunc(w io.Writer, format string, a ...interface{})

	// Writer gets underlying io.Writer
	Writer() io.Writer
}

// Printer represents printing abstraction
type Printer interface {

	// Writer gets underlying io.Writer
	Writer() io.Writer

	// Twriter gets underlying  *tabwriter.Writer
	// to output tabular data
	Twriter() *tabwriter.Writer

	// Flush flushes tabular writer that prints dato
	Flush()

	// Tprint prints using tab writer
	Tprint(format string, a ...interface{})

	// Cprint prints data with suppport colorizing
	Cprint(format string, a ...interface{})
}

// Executor represent executable command interface
type Executor interface {
	// Execute starts execute command's code
	Execute() error
}
