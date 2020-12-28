package cmd

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
)

type validateCommand struct{ baseCommand }
type fixCommand struct{ baseCommand }

func newValidate(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() api.Executor {
			return &validateCommand{
				baseCommand: newBaseCmd(c),
			}
		},
		c: c,
	}

	cmd := cc.NewCommand("va", "validate", "Validates SDK projects within solution(s)")

	cmd.AddCommand(newFix(c))

	return cmd
}

func newFix(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() api.Executor {
			return &fixCommand{
				baseCommand: newBaseCmd(c),
			}
		},
		c: c,
	}

	cmd := cc.NewCommand("fix", "fixprojects", "Fixes redundant SDK projects references")

	return cmd
}

func (c *validateCommand) Execute() error {
	projectsPrinter := newSdkProjectsPrinter(c.prn)
	validator := newSdkProjectsValidator(c.fs, c.prn, c.sourcesPath, projectsPrinter)

	validator.validate()
	return nil
}

func (c *fixCommand) Execute() error {
	fixer := newSdkProjectsFixer(c.prn, c.fs)
	validator := newSdkProjectsValidator(c.fs, c.prn, c.sourcesPath, fixer)

	validator.validate()
	return nil
}
