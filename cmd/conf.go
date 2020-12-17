package cmd

import (
	"github.com/spf13/afero"
)

type appConf struct {
	filesystem afero.Fs
	p          printer
	sp         *string
}

func (a *appConf) fs() afero.Fs {
	return a.filesystem
}

func (a *appConf) prn() printer {
	return a.p
}

func (a *appConf) sourcesPath() *string {
	return a.sp
}

func newAppConf(fs afero.Fs, p printer, sourcesPath *string) conf {
	c := appConf{
		filesystem: fs,
		p:          p,
		sp:         sourcesPath,
	}
	return &c
}
