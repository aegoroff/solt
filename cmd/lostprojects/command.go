package lostprojects

import (
	"github.com/spf13/cobra"
	"solt/cmd/fw"
	"solt/msvc"
)

type lostProjectsCommand struct {
	*fw.BaseCommand
}

// New creates new command that does lost projects search
func New(c *fw.Conf) *cobra.Command {
	cc := fw.NewCobraCreator(c, func() fw.Executor {
		exe := &lostProjectsCommand{fw.NewBaseCmd(c)}
		return fw.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("lp", "lostprojects", "Find projects that not included into any solution")

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

	lost, lostWithIncludes := findLostProjects(allProjects, solutionProjects.Includes())

	s := fw.NewScreener(c.Prn())
	// Lost projects
	s.WriteSlice(lost)

	if len(lostWithIncludes) > 0 {
		m := "\n<red>These projects are not included into any solution but files from the projects' folders are used in another projects within a solution:</>\n\n"
		c.Prn().Cprint(m)
	}

	// Lost projects that have includes files that used
	s.WriteSlice(lostWithIncludes)

	title := "<red>These projects are included into a solution but not found in the file system:</>"
	exist.Print(c.Prn(), title, "Solution")

	return nil
}

func findLostProjects(allProjects []*msvc.MsbuildProject, solutionProjects []string) ([]string, []string) {
	// Create projects matching machine
	incl := fw.NewExactMatch(solutionProjects)

	lostProjects := allProjects[:0]

	allSolutionFiles := fw.NewIncluder(fw.NewNullExister())

	for _, prj := range allProjects {
		if incl.Match(prj.Path()) {
			allSolutionFiles.From(prj)
		} else {
			lostProjects = append(lostProjects, prj)
		}
	}

	return separateProjects(lostProjects, allSolutionFiles.Includes())
}

func separateProjects(lostProjects []*msvc.MsbuildProject, allSolutionFiles []string) ([]string, []string) {
	var lost []string
	var lostWithIncludes []string
	solutionFiles := fw.NewExactMatch(allSolutionFiles)

	for _, lp := range lostProjects {
		if fw.MatchAny(lp.Items(), solutionFiles) {
			lostWithIncludes = append(lostWithIncludes, lp.Path())
		} else {
			lost = append(lost, lp.Path())
		}
	}
	return lost, lostWithIncludes
}
