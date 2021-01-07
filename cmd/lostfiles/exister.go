package lostfiles

import (
	"github.com/spf13/afero"
	"io"
	"solt/cmd/fw"
)

func newExister(validate bool, fs afero.Fs, w io.Writer) fw.Exister {
	if validate {
		return fw.NewExister(fs, w)
	}
	return &stubExister{}
}

// Stub that do nothing

type stubExister struct{}

func (*stubExister) Print(fw.Printer, string, string) {}

func (*stubExister) Validate(string, []string) {}
