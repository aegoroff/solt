package api

import (
	"github.com/spf13/afero"
)

type Conf struct {
	filesystem afero.Fs
	p          Printer
	sp         *string
	cpu        *string
	memory     *string
	diag       *bool
}

func NewConf(fs afero.Fs, p Printer, sp *string, cpu *string, memory *string, diag *bool) *Conf {
	return &Conf{filesystem: fs, p: p, sp: sp, cpu: cpu, memory: memory, diag: diag}
}

func (a *Conf) Diag() *bool {
	return a.diag
}

func (a *Conf) Cpu() *string {
	return a.cpu
}

func (a *Conf) Memory() *string {
	return a.memory
}

func (a *Conf) Fs() afero.Fs {
	return a.filesystem
}

func (a *Conf) Prn() Printer {
	return a.p
}

func (a *Conf) SourcesPath() *string {
	return a.sp
}
