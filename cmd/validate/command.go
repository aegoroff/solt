package validate

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
)

type validateCommand struct{ api.BaseCommand }
type fixCommand struct{ api.BaseCommand }

func New(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		return &validateCommand{
			BaseCommand: api.NewBaseCmd(c),
		}
	})

	cmd := cc.NewCommand("va", "validate", "Validates SDK projects within solution(s)")

	cmd.AddCommand(newFix(c))

	return cmd
}

func newFix(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		return &fixCommand{
			BaseCommand: api.NewBaseCmd(c),
		}
	})

	cmd := cc.NewCommand("fix", "fixprojects", "Fixes redundant SDK projects references")

	return cmd
}

func (c *validateCommand) Execute() error {
	projectsPrinter := newSdkProjectsPrinter(c.Prn())
	validator := newSdkProjectsValidator(c.Fs(), c.Prn(), c.SourcesPath(), projectsPrinter)

	validator.validate()
	return nil
}

func (c *fixCommand) Execute() error {
	fixer := newSdkProjectsFixer(c.Prn(), c.Fs())
	validator := newSdkProjectsValidator(c.Fs(), c.Prn(), c.SourcesPath(), fixer)

	validator.validate()
	return nil
}
