package cmd

import (
	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.11.0"

type versionCommand struct {
	baseCommand
}

func newVersion(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() executor {
			vac := versionCommand{
				baseCommand: newBaseCmd(c),
			}
			return &vac
		},
		c: c,
	}

	cmd := cc.newCobraCommand("ver", "version", "Print the version number of solt")
	cmd.Long = `All software has versions. This is solt's`

	return cmd
}

func (c *versionCommand) execute() error {
	c.prn.cprint("%s\n", Version)
	return nil
}
