package cmd

import (
	"github.com/spf13/cobra"
)

type validateCommand struct {
	baseCommand
}

func newValidate(c conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() command {
			return &validateCommand{
				baseCommand: newBaseCmd(c),
			}
		},
	}

	cmd := cc.newCobraCommand("va", "validate", "Validates SDK projects within solution(s)")

	return cmd
}

func (c *validateCommand) execute() error {
	v := newsdkProjectsValidator(c.prn)
	m := newSdkProjectsModule(c.fs, c.prn, c.sourcesPath, v)

	m.execute()
	return nil
}
