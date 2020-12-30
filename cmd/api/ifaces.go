package api

import (
	"fmt"
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
	Writable

	// PrintFunc represents printing function implementation
	PrintFunc(w io.Writer, format string, a ...interface{})

	// NewPrinter creates new printer
	NewPrinter() Printer
}

// StringEnvironment defines in memory printing environment abstraction
type StringEnvironment interface {
	PrintEnvironment
	fmt.Stringer
}

// Printer represents printing abstraction
type Printer interface {
	Writable

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

// Writable represents io.Writer container
type Writable interface {
	// Writer gets underlying io.Writer
	Writer() io.WriteCloser
}

// Executor represent executable command interface
type Executor interface {
	// Execute starts execute command's code
	Execute() error
}
