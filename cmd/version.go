package cmd

import (
	"github.com/spf13/cobra"
	"solt/cmd/fw"
)

// Version defines program version
var Version = "0.14.0-dev"

type versionCommand struct {
	*fw.BaseCommand
}

func newVersion(c *fw.Conf) *cobra.Command {
	cc := fw.NewCobraCreator(c, func() fw.Executor {
		return &versionCommand{
			BaseCommand: fw.NewBaseCmd(c),
		}
	})

	cmd := cc.NewArgsCommand("ver", "version", "Print the version number of solt", nil)
	cmd.Long = `All software has versions. This is solt's`

	return cmd
}

func (c *versionCommand) Execute(*cobra.Command) error {
	c.Prn().Cprint("%s\n", Version)
	return nil
}
