package lostprojects

import (
	"github.com/spf13/afero"
	"io"
	"solt/cmd/api"
	"solt/internal/sys"
)

type exister struct {
	filer   sys.Filer
	unexist map[string][]string
}

func newExister(fs afero.Fs, w io.Writer) *exister {
	return &exister{
		unexist: make(map[string][]string),
		filer:   sys.NewFiler(fs, w),
	}
}

func (e *exister) validate(sol string, paths []string) {
	nonexist := e.filer.CheckExistence(paths)

	if len(nonexist) > 0 {
		e.unexist[sol] = append(e.unexist[sol], nonexist...)
	}
}

func (e *exister) print(p api.Printer) {
	if len(e.unexist) > 0 {
		p.Cprint("\n<red>These projects are included into a solution but not found in the file system:</>\n")
	}

	s := api.NewScreener(p)
	s.WriteMap(e.unexist, "Solution")
}
