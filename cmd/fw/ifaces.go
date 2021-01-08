package fw

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
)

// Matcher defines string matcher interface
type Matcher interface {
	// Match does string matching to several patterns
	Match(s string) bool
}

// Searcher defines string searching interface
type Searcher interface {
	Matcher
	// Search does string matching to several patterns and returns all found
	Search(s string) []string
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

// Executor represent executable command interface
type Executor interface {
	// Execute starts execute command's code
	Execute(cc *cobra.Command) error
}

// Exister provides file existence validation in a container
type Exister interface {
	// Validate validates whether files from container exist in filesystem
	Validate(root string, includes []string)
	// Print outputs unexist files info
	Print(p Printer, title string, container string)
}
