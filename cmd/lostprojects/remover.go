package lostprojects

import (
	"github.com/spf13/afero"
	"solt/cmd/out"
)

type remover struct {
	fs     afero.Fs
	remove bool
	p      out.Printer
}

func newRemover(fs afero.Fs, p out.Printer, remove bool) *remover {
	return &remover{fs: fs, remove: remove, p: p}
}

func (r *remover) removeAll(projects []string) error {
	if !r.remove {
		return nil
	}
	for _, p := range projects {
		d := dir(p)
		err := r.fs.RemoveAll(d)
		if err != nil {
			return err
		}
		r.p.Cprint("\n Folder '<red>%s</>' removed", d)
	}
	if len(projects) > 0 {
		r.p.Println()
	}
	return nil
}
