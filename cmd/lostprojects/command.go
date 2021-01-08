package lostprojects

import (
	"github.com/spf13/cobra"
	"solt/cmd/fw"
	"solt/msvc"
)

type lostProjectsCommand struct {
	*fw.BaseCommand
	removeLost bool
}

// New creates new command that does lost projects search
func New(c *fw.Conf) *cobra.Command {
	var removeLost bool

	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &lostProjectsCommand{
			BaseCommand: fw.NewBaseCmd(c),
			removeLost:  removeLost,
		}
		return fw.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("lp", "lostprojects", "Find projects that not included into any solution")
	cmd.Flags().BoolVarP(&removeLost, "remove", "r", false, "Remove lost projects folders")

	return cmd
}

func (c *lostProjectsCommand) Execute(*cobra.Command) error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())

	solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

	exist := fw.NewExister(c.Fs(), c.Writer())
	solutionProjects := fw.NewIncluder(exist)

	// Each found solution
	for _, sln := range solutions {
		solutionProjects.From(sln)
	}

	find := newFinder()

	lost, lostWithIncludes := find.filter(allProjects, solutionProjects.Includes())

	s := fw.NewScreener(c.Prn())
	// Lost projects
	s.WriteSlice(lost)

	if len(lostWithIncludes) > 0 {
		m1 := fw.NewMarginer(1)

		l1 := "<red>These projects are not included into any solution</>"
		l2 := "<red>but files from the projects' folders are used in another projects within a solution:</>"
		c.Prn().Println()
		c.Prn().Cprint(m1.Margin(l1))
		c.Prn().Println()
		c.Prn().Cprint(m1.Margin(l2))
		c.Prn().Println()
		c.Prn().Println()

		s.WriteSlice(lostWithIncludes)
	}

	title := "<red>These projects are included into a solution but not found in the file system:</>"
	exist.Print(c.Prn(), title, "Solution")

	r := newRemover(c.Fs(), c.Prn(), c.removeLost)
	return r.removeAll(lost)
}
