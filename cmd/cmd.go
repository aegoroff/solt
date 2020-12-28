package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/cmd/api"
)

type cobraRunSignature func(cmd *cobra.Command, args []string) error

type baseCommand struct {
	prn         api.Printer
	fs          afero.Fs
	sourcesPath string
}

func newBaseCmd(c *conf) baseCommand {
	return baseCommand{
		prn:         c.prn(),
		sourcesPath: *c.sourcesPath(),
		fs:          c.fs(),
	}
}

type cobraCreator struct {
	createCmd func() api.Executor
	c         *conf
}

func (c *cobraCreator) runE() cobraRunSignature {
	return func(cmd *cobra.Command, args []string) error {
		// IMPORTANT: Excecutors initialization order defines output order
		var e api.Executor
		{
			e = c.createCmd()
			e = newCPUProfileExecutor(e, c.c)
			e = newMemUsageExecutor(e, c.c)
			e = newTimeMeasureExecutor(e, c.c)
			e = newMemoryProfileExecutor(e, c.c)
		}

		return e.Execute()
	}
}

func (c *cobraCreator) NewCommand(use, alias, short string) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     use,
		Aliases: []string{alias},
		Short:   short,
		RunE:    c.runE(),
	}
	return cmd
}
