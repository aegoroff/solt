package fw

import (
	"github.com/spf13/cobra"
	"solt/internal/out"
	"solt/msvc"
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
	Print(p out.Printer, title string, container string)
	// UnexistCount gets the number of non exist items
	UnexistCount() int64
}

// Solutioner provides solution action prototype
type Solutioner interface {
	// Solution method called on each solution
	Solution(*msvc.VisualStudioSolution)
}
