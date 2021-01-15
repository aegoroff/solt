package info

import (
	"github.com/spf13/cobra"
	"solt/cmd/fw"
	"solt/msvc"
)

type infoCommand struct {
	*fw.BaseCommand
}

// New creates new command that shows information about solutions
func New(c *fw.Conf) *cobra.Command {
	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &infoCommand{
			BaseCommand: fw.NewBaseCmd(c),
		}
		return fw.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("in", "info", "Get information about found solutions")
	return cmd
}

func (c *infoCommand) Execute(*cobra.Command) error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())

	solutions := msvc.SelectSolutions(foldersTree)
	msvc.SortSolutions(solutions)

	handlers := []solutioner{
		newDisplay(c.Prn(), c),
	}

	for _, sol := range solutions {
		for _, h := range handlers {
			h.solution(sol)
		}
	}

	return nil
}
