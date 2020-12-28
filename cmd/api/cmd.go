package api

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

type cobraRunSignature func(cmd *cobra.Command, args []string) error

type BaseCommand struct {
	prn         Printer
	fs          afero.Fs
	sourcesPath string
}

func (b *BaseCommand) SourcesPath() string {
	return b.sourcesPath
}

func (b *BaseCommand) Fs() afero.Fs {
	return b.fs
}

func (b *BaseCommand) Prn() Printer {
	return b.prn
}

func NewBaseCmd(c *Conf) BaseCommand {
	return BaseCommand{
		prn:         c.Prn(),
		sourcesPath: *c.SourcesPath(),
		fs:          c.Fs(),
	}
}

type CobraCreator struct {
	exe func() Executor
	c   *Conf
}

func NewCobraCreator(c *Conf, exe func() Executor) *CobraCreator {
	return &CobraCreator{exe: exe, c: c}
}

func (c *CobraCreator) runE() cobraRunSignature {
	return func(cc *cobra.Command, args []string) error {
		// IMPORTANT: Excecutors initialization order defines output order
		var e Executor
		{
			e = c.exe()
			e = newCPUProfileExecutor(e, c.c)
			e = newMemUsageExecutor(e, c.c)
			e = newTimeMeasureExecutor(e, c.c)
			e = newMemoryProfileExecutor(e, c.c)
		}

		return e.Execute()
	}
}

func (c *CobraCreator) NewCommand(use, alias, short string) *cobra.Command {
	var cc = &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   short,
		RunE:    c.runE(),
	}
	return cc
}
