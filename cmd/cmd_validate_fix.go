package cmd

import "github.com/spf13/cobra"

type fixCommand struct {
	baseCommand
}

func newFix(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() executor {
			return &fixCommand{
				baseCommand: newBaseCmd(c),
			}
		},
	}

	cmd := cc.newCobraCommand("fix", "fixprojects", "Fixes redundant SDK projects references")

	return cmd
}

func (c *fixCommand) execute() error {
	fixer := newSdkProjectsFixer(c.prn, c.fs)
	validator := newSdkProjectsValidator(c.fs, c.prn, c.sourcesPath, fixer)

	validator.validate()
	return nil
}
