package fw

import (
	"github.com/spf13/afero"
	"io"
	"solt/internal/out"
	"solt/internal/sys"
	"solt/internal/ux"
)

type exister struct {
	filer        Filer
	unexist      map[string][]string
	unexistCount int64
}

// NewExister creates new Exister instance
func NewExister(fs afero.Fs, w io.Writer) Exister {
	return &exister{
		unexist: make(map[string][]string),
		filer:   sys.NewFiler(fs, w),
	}
}

// Validate validates whether files from container exist in filesystem
func (e *exister) Validate(root string, paths []string) {
	nonexist := e.filer.CheckExistence(paths)
	l := len(nonexist)
	e.unexistCount += int64(l)

	if l > 0 {
		e.unexist[root] = append(e.unexist[root], nonexist...)
	}
}

// UnexistCount gets the number of non exist items
func (e *exister) UnexistCount() int64 {
	return e.unexistCount
}

// Print outputs unexist files in
func (e *exister) Print(p out.Printer, title string, container string) {
	if len(e.unexist) > 0 {
		p.Println()
		p.Cprint(title)
		p.Println()
	}

	s := ux.NewScreener(p)
	s.WriteMap(e.unexist, container)
}
