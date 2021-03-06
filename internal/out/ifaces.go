package out

import (
	"fmt"
	"io"
)

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

// Printer represents printing abstraction with colorizing support
type Printer interface {
	// Cprint prints data with colorizing support
	Cprint(format string, a ...interface{})

	// Println prints new line
	Println()
}

// Writable represents io.Writer container
type Writable interface {
	// Writer gets underlying io.Writer
	Writer() io.WriteCloser
}
