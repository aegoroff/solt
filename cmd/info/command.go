package info

import (
	"github.com/spf13/cobra"
	"solt/internal/fw"
	"solt/msvc"
	"sort"
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

	sols := msvc.SelectSolutions(foldersTree)

	solutions := fw.SolutionSlice(sols)
	sort.Sort(solutions)

	grp := newProjectGroupper()
	collect := newCollector(grp)
	solutions.Foreach(grp, newDisplay(c.Prn(), c, grp), collect)

	c.Total(collect.result)

	return nil
}
