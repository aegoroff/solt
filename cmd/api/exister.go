package api

import (
	"github.com/spf13/afero"
	"io"
	"solt/internal/sys"
)

// Exister provides file existence validation in a container
type Exister struct {
	filer   sys.Filer
	unexist map[string][]string
}

// NewExister creates new Exister instance
func NewExister(fs afero.Fs, w io.Writer) *Exister {
	return &Exister{
		unexist: make(map[string][]string),
		filer:   sys.NewFiler(fs, w),
	}
}

// Validate validates whether files from container exist in filesystem
func (e *Exister) Validate(root string, paths []string) {
	nonexist := e.filer.CheckExistence(paths)

	if len(nonexist) > 0 {
		e.unexist[root] = append(e.unexist[root], nonexist...)
	}
}

// Print outputs unexist files info
func (e *Exister) Print(p Printer, title string, container string) {
	if len(e.unexist) > 0 {
		p.Println()
		p.Cprint(title)
		p.Println()
	}

	s := NewScreener(p)
	s.WriteMap(e.unexist, container)
}
