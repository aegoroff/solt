package cmd

import (
	"github.com/spf13/cobra"
)

type validateCommand struct {
	baseCommand
}

func newValidate(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() executor {
			return &validateCommand{
				baseCommand: newBaseCmd(c),
			}
		},
	}

	cmd := cc.newCobraCommand("va", "validate", "Validates SDK projects within solution(s)")

	cmd.AddCommand(newFix(c))

	return cmd
}

func (c *validateCommand) execute() error {
	projectsPrinter := newSdkProjectsPrinter(c.prn)
	validator := newSdkProjectsValidator(c.fs, c.prn, c.sourcesPath, projectsPrinter)

	validator.validate()
	return nil
}
