package cmd

import (
	"github.com/spf13/cobra"
)

// Version defines program version
var Version = "0.9.2"

type versionCommand struct {
	baseCommand
}

func newVersion(c conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() command {
			vac := versionCommand{
				baseCommand: newBaseCmd(c),
			}
			return &vac
		},
	}

	cmd := cc.newCobraCommand("ver", "version", "Print the version number of solt")
	cmd.Long = `All software has versions. This is solt's`

	return cmd
}

func (c *versionCommand) execute() error {
	c.prn.cprint("solt v%s\n", Version)
	return nil
}
