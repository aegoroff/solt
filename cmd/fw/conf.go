package fw

import (
	"github.com/spf13/afero"
	"io"
	"solt/internal/out"
)

// Conf is app configuration container
type Conf struct {
	filesystem afero.Fs
	p          out.Printer
	pe         out.PrintEnvironment
	sp         *string
	cpu        *string
	memory     *string
	diag       *bool
}

// Diag is app diagnostic context
type Diag struct {
	// CPU path to cpu profiling results
	CPU string
	// Memory path to memory profiling results
	Memory string
	// Enable whether to enable diagnostic
	Enable bool
}

// NewConf creates new *Conf instance
func NewConf(fs afero.Fs, pe out.PrintEnvironment, d *Diag) *Conf {
	return &Conf{filesystem: fs, pe: pe, cpu: &d.CPU, memory: &d.Memory, diag: &d.Enable}
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
func (c *Conf) Prn() out.Printer {
	return c.p
}

// W gets underlying Writer
func (c *Conf) W() io.WriteCloser {
	return c.pe.Writer()
}

// SourcesPath gets analyzable sources path
func (c *Conf) SourcesPath() *string {
	return c.sp
}

func (c *Conf) init(sources *string) error {
	p, err := c.pe.NewPrinter()
	c.p = p
	c.sp = sources
	if err != nil || *c.SourcesPath() == "" {
		return err
	}

	sp := *c.SourcesPath()
	_, err = c.Fs().Stat(sp)

	return err
}
