package validate

import (
	"github.com/spf13/cobra"
	"solt/cmd/fw"
)

type validateCommand struct{ *fw.BaseCommand }
type fixCommand struct{ *fw.BaseCommand }

// New creates new command that does validates SDK projects
func New(c *fw.Conf) *cobra.Command {
	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &validateCommand{fw.NewBaseCmd(c)}
		return fw.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("va", "validate", "Validates SDK projects within solution(s)")

	cmd.AddCommand(newFix(c))

	return cmd
}

func newFix(c *fw.Conf) *cobra.Command {
	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &fixCommand{fw.NewBaseCmd(c)}
		return fw.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("fix", "fixprojects", "Fixes redundant SDK projects references")

	return cmd
}

func (c *validateCommand) Execute(*cobra.Command) error {
	prn := newPrinter(c.Prn())
	valid := newValidator(c.Fs(), c.SourcesPath(), prn)

	valid.validate()

	valid.tt.display(c.Prn(), c)
	return nil
}

func (c *fixCommand) Execute(*cobra.Command) error {
	fix := newFixer(c.Prn(), c, c.Fs())
	valid := newValidator(c.Fs(), c.SourcesPath(), fix)

	valid.validate()
	return nil
}
