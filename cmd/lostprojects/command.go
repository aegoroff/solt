package lostprojects

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/msvc"
)

type lostProjectsCommand struct {
	*api.BaseCommand
}

// New creates new command that does lost projects search
func New(c *api.Conf) *cobra.Command {
	cc := api.NewCobraCreator(c, func() api.Executor {
		exe := &lostProjectsCommand{api.NewBaseCmd(c)}
		return api.NewExecutorShowHelp(exe, c)
	})

	cmd := cc.NewCommand("lp", "lostprojects", "Find projects that not included into any solution")

	return cmd
}

func (c *lostProjectsCommand) Execute(*cobra.Command) error {
	foldersTree := msvc.ReadSolutionDir(c.SourcesPath(), c.Fs())

	solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

	// linked from any solution projects list
	// so these projects are not considered lost
	var linkedProjects []string

	exist := api.NewExister(c.Fs(), c.Writer())

	// Each found solution
	for _, sln := range solutions {
		projects := sln.AllProjectPaths(msvc.PassThrough)
		linkedProjects = append(linkedProjects, projects...)
		exist.Validate(sln.Path, projects)
	}

	lost, lostWithIncludes := findLostProjects(allProjects, linkedProjects)

	s := api.NewScreener(c.Prn())
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

func findLostProjects(allProjects []*msvc.MsbuildProject, linkedProjects []string) ([]string, []string) {
	// Create projects matching machine
	incl := api.NewExactMatch(linkedProjects)

	lostProjects := allProjects[:0]
	var allSolutionFiles []string

	for _, prj := range allProjects {
		if !incl.Match(prj.Path) {
			lostProjects = append(lostProjects, prj)
		} else {
			allSolutionFiles = append(allSolutionFiles, prj.Files()...)
		}
	}

	return separateProjects(lostProjects, allSolutionFiles)
}

func separateProjects(projectsOutsideSolution []*msvc.MsbuildProject, allSolutionFiles []string) ([]string, []string) {
	var lost []string
	var lostWithIncludes []string
	solutionFiles := api.NewExactMatch(allSolutionFiles)

	for _, prj := range projectsOutsideSolution {
		if solutionFiles.MatchAny(prj.Files()) {
			lostWithIncludes = append(lostWithIncludes, prj.Path)
		} else {
			lost = append(lost, prj.Path)
		}
	}
	return lost, lostWithIncludes
}
