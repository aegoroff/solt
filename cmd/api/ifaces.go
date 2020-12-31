package api

import (
	"fmt"
	"io"
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
	NewPrinter() (Printer, error)
}

// StringEnvironment defines in memory printing environment abstraction
type StringEnvironment interface {
	PrintEnvironment
	fmt.Stringer
}

// Printer represents printing abstraction
type Printer interface {
	Writable
	ColorPrinter
}

// ColorPrinter represents printing abstraction with colorizing support
type ColorPrinter interface {
	// Cprint prints data with colorizing support
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
