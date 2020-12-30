package api

import (
	"github.com/spf13/afero"
)

// Conf is app configuration container
type Conf struct {
	filesystem afero.Fs
	p          Printer
	pe         PrintEnvironment
	sp         *string
	cpu        *string
	memory     *string
	diag       *bool
}

// NewConf creates new *Conf instance
func NewConf(fs afero.Fs, pe PrintEnvironment, sp *string, cpu *string, memory *string, diag *bool) *Conf {
	return &Conf{filesystem: fs, pe: pe, sp: sp, cpu: cpu, memory: memory, diag: diag}
}

// Diag gets whether to enable diagnostic mode
func (c *Conf) Diag() *bool {
	return c.diag
}

// CPU gets cpu profiling file path that will be created
func (c *Conf) CPU() *string {
	return c.cpu
}

// Memory gets memory profiling file path that will be created
func (c *Conf) Memory() *string {
	return c.memory
}

// Fs gets underlying file system abstraction
func (c *Conf) Fs() afero.Fs {
	return c.filesystem
}

// Prn gets underlying Printer
func (c *Conf) Prn() Printer {
	return c.p
}

// SourcesPath gets analyzable sources path
func (c *Conf) SourcesPath() *string {
	return c.sp
}

func (c *Conf) init() {
	c.p = c.pe.NewPrinter()
}
