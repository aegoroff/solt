package lostfiles

import (
	"github.com/spf13/afero"
	"io"
	"solt/cmd/fw"
)

func newExister(validate bool, fs afero.Fs, w io.Writer) exister {
	if validate {
		return &realExister{e: fw.NewExister(fs, w)}
	}
	return &stubExister{}
}

// Stub that do nothing

type stubExister struct{}

func (*stubExister) print(fw.Printer) {}

func (*stubExister) exist(string, []string) {}

// Real

type realExister struct {
	e *fw.Exister
}

func (r *realExister) print(p fw.Printer) {
	title := "<red>These files included into projects but not exist in the file system.</>"
	r.e.Print(p, title, "Project")
}

func (r *realExister) exist(project string, includes []string) {
	r.e.Validate(project, includes)
}
