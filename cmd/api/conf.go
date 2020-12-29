package api

import (
	"github.com/spf13/afero"
)

// Conf is app configuration container
type Conf struct {
	filesystem afero.Fs
	p          Printer
	sp         *string
	cpu        *string
	memory     *string
	diag       *bool
}

// NewConf creates new *Conf instance
func NewConf(fs afero.Fs, p Printer, sp *string, cpu *string, memory *string, diag *bool) *Conf {
	return &Conf{filesystem: fs, p: p, sp: sp, cpu: cpu, memory: memory, diag: diag}
}

// Diag gets whether to enable diagnostic mode
func (a *Conf) Diag() *bool {
	return a.diag
}

// CPU gets cpu profiling file path that will be created
func (a *Conf) CPU() *string {
	return a.cpu
}

// Memory gets memory profiling file path that will be created
func (a *Conf) Memory() *string {
	return a.memory
}

// Fs gets underlying file system abstraction
func (a *Conf) Fs() afero.Fs {
	return a.filesystem
}

// Prn gets underlying Printer
func (a *Conf) Prn() Printer {
	return a.p
}

// SourcesPath gets analyzable sources path
func (a *Conf) SourcesPath() *string {
	return a.sp
}
