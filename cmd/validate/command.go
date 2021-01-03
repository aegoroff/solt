package validate

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
)

type validateCommand struct{ *api.BaseCommand }
type fixCommand struct{ *api.BaseCommand }

// New creates new command that does validates SDK projects
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

func (c *validateCommand) Execute(cc *cobra.Command) error {
	prn := newPrinter(c.Prn())
	valid := newValidator(c.Fs(), c.SourcesPath(), prn)

	valid.validate()
	return c.ShowHelp(cc)
}

func (c *fixCommand) Execute(cc *cobra.Command) error {
	fix := newFixer(c.Prn(), c, c.Fs())
	valid := newValidator(c.Fs(), c.SourcesPath(), fix)

	valid.validate()
	return c.ShowHelp(cc)
}
