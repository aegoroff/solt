package cmd

import (
	"github.com/spf13/afero"
)

type conf struct {
	filesystem afero.Fs
	p          printer
	sp         *string
	cpu        *string
	memory     *string
	diag       *bool
}

func (a *conf) fs() afero.Fs {
	return a.filesystem
}

func (a *conf) prn() printer {
	return a.p
}

func (a *conf) sourcesPath() *string {
	return a.sp
}
