package cmd

import (
	"github.com/spf13/afero"
	"solt/cmd/api"
)

type conf struct {
	filesystem afero.Fs
	p          api.Printer
	sp         *string
	cpu        *string
	memory     *string
	diag       *bool
}

func (a *conf) fs() afero.Fs {
	return a.filesystem
}

func (a *conf) prn() api.Printer {
	return a.p
}

func (a *conf) sourcesPath() *string {
	return a.sp
}
