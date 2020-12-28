package cmd

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
)

// Version defines program version
var Version = "0.11.0"

type versionCommand struct {
	api.BaseCommand
}

func newVersion(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		return &versionCommand{
			BaseCommand: api.NewBaseCmd(c),
		}
	})

	cmd := cc.NewCommand("ver", "version", "Print the version number of solt")
	cmd.Long = `All software has versions. This is solt's`

	return cmd
}

func (c *versionCommand) Execute() error {
	c.Prn().Cprint("%s\n", Version)
	return nil
}
