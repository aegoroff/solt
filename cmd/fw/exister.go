package fw

import (
	"github.com/spf13/afero"
	"io"
	"solt/internal/sys"
)

type exister struct {
	filer   sys.Filer
	unexist map[string][]string
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

	if len(nonexist) > 0 {
		e.unexist[root] = append(e.unexist[root], nonexist...)
	}
}

// Print outputs unexist files info
func (e *exister) Print(p Printer, title string, container string) {
	if len(e.unexist) > 0 {
		p.Println()
		p.Cprint(title)
		p.Println()
	}

	s := NewScreener(p)
	s.WriteMap(e.unexist, container)
}

// NewNullExister creates new Exister that do nothing
func NewNullExister() Exister { return &nullExister{} }

type nullExister struct{}

func (*nullExister) Print(Printer, string, string) {}

func (*nullExister) Validate(string, []string) {}
