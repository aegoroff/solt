package cmd

import (
	"github.com/spf13/cobra"
)

type validateCommand struct{ baseCommand }
type fixCommand struct{ baseCommand }

func newValidate(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() executor {
			return &validateCommand{
				baseCommand: newBaseCmd(c),
			}
		},
		c: c,
	}

	cmd := cc.newCobraCommand("va", "validate", "Validates SDK projects within solution(s)")

	cmd.AddCommand(newFix(c))

	return cmd
}

func newFix(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() executor {
			return &fixCommand{
				baseCommand: newBaseCmd(c),
			}
		},
		c: c,
	}

	cmd := cc.newCobraCommand("fix", "fixprojects", "Fixes redundant SDK projects references")

	return cmd
}

func (c *validateCommand) execute() error {
	projectsPrinter := newSdkProjectsPrinter(c.prn)
	validator := newSdkProjectsValidator(c.fs, c.prn, c.sourcesPath, projectsPrinter)

	validator.validate()
	return nil
}

func (c *fixCommand) execute() error {
	fixer := newSdkProjectsFixer(c.prn, c.fs)
	validator := newSdkProjectsValidator(c.fs, c.prn, c.sourcesPath, fixer)

	validator.validate()
	return nil
}
