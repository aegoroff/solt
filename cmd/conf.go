package cmd

import (
	"github.com/spf13/afero"
	"io"
)

type conf interface {
	// fs defines app file system abstraction
	fs() afero.Fs

	prn() printer

	globals() *globals
}

type appConf struct {
	filesystem afero.Fs
	p          printer
	g          *globals
}

type globals struct {
	sourcesPath *string
	diag        *bool
}

func (a *appConf) fs() afero.Fs {
	return a.filesystem
}

func (a *appConf) prn() printer {
	return a.p
}

func (a *appConf) globals() *globals {
	return a.g
}

func newAppConf(fs afero.Fs, w io.Writer, g *globals) conf {
	c := appConf{
		filesystem: fs,
		p:          newPrinter(w),
		g:          g,
	}
	return &c
}
