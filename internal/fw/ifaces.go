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

// Solutioner provides solution action prototype
type Solutioner interface {
	// Solution method called on each solution
	Solution(*msvc.VisualStudioSolution)
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
	// Print outputs unexist files in
	Print(p out.Printer, title string, container string)
	// MissingCount gets the number of non exist items
	MissingCount() int64
}

// Filer defines module that works with files
type Filer interface {
	Remover
	// CheckExistence validates files passed to be present in file system
	// The list of non exist files returned
	CheckExistence(files []string) []string

	// Write writes new file content
	Write(path string, content []byte)

	// Read reads file content
	Read(path string) ([]byte, error)
}

// Remover defines removing files interface
type Remover interface {
	// Remove removes files from file system
	Remove(files []string)
}

// container provides paths container interface
type container interface {
	// Path provides container's path itself
	Path() string
	// Items provides all paths included into container
	Items() []string
}
