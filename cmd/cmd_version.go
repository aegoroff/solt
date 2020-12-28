package cmd

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
)

// Version defines program version
var Version = "0.11.0"

type versionCommand struct {
	baseCommand
}

func newVersion(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() api.Executor {
			vac := versionCommand{
				baseCommand: newBaseCmd(c),
			}
			return &vac
		},
		c: c,
	}

	cmd := cc.NewCobraCommand("ver", "version", "Print the version number of solt")
	cmd.Long = `All software has versions. This is solt's`

	return cmd
}

func (c *versionCommand) Execute() error {
	c.prn.Cprint("%s\n", Version)
	return nil
}
