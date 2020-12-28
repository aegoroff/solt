package cmd

import (
	"github.com/spf13/cobra"
	"solt/cmd/api"
	"solt/internal/sys"
	"solt/msvc"
)

type lostProjectsCommand struct {
	baseCommand
}

func newLostProjects(c *conf) *cobra.Command {
	cc := cobraCreator{
		createCmd: func() api.Executor {
			lpc := lostProjectsCommand{
				baseCommand: newBaseCmd(c),
			}
			return &lpc
		},
		c: c,
	}

	cmd := cc.NewCommand("lp", "lostprojects", "Find projects that not included into any solution")

	return cmd
}

func (c *lostProjectsCommand) Execute() error {
	foldersTree := msvc.ReadSolutionDir(c.sourcesPath, c.fs)

	solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

	// linked from any solution projects list
	// so these projects are not considered lost
	var linkedProjects []string

	projectLinksBySolution := make(map[string][]string)
	// Each found solution
	for _, sln := range solutions {
		links := sln.AllProjectPaths(msvc.PassThrough)
		projectLinksBySolution[sln.Path] = links
		linkedProjects = append(linkedProjects, links...)
	}

	lost, lostWithIncludes := findLostProjects(allProjects, linkedProjects)

	s := newScreener(c.prn)
	// Lost projects
	s.writeSlice(lost)

	if len(lostWithIncludes) > 0 {
		m := "\n<red>These projects are not included into any solution but files from the projects' folders are used in another projects within a solution:</>\n\n"
		c.prn.Cprint(m)
	}

	// Lost projects that have includes files that used
	s.writeSlice(lostWithIncludes)

	unexistProjects := c.getUnexistProjects(projectLinksBySolution)

	if len(unexistProjects) > 0 {
		c.prn.Cprint("\n<red>These projects are included into a solution but not found in the file system:</>\n")
	}

	// Included but not exist in FS
	s.writeMap(unexistProjects, "Solution")

	return nil
}

func (c *lostProjectsCommand) getUnexistProjects(projectsInSolutions map[string][]string) map[string][]string {
	var result = make(map[string][]string)

	filer := sys.NewFiler(c.fs, c.prn.Writer())
	for spath, projects := range projectsInSolutions {
		nonexist := filer.CheckExistence(projects)

		if len(nonexist) > 0 {
			result[spath] = append(result[spath], nonexist...)
		}
	}
	return result
}

func findLostProjects(allProjects []*msvc.MsbuildProject, linkedProjects []string) ([]string, []string) {
	// Create projects matching machine
	projectMatch := NewExactMatch(linkedProjects)
	projectsOutsideSolution := allProjects[:0]
	var allSolutionFiles []string

	for _, prj := range allProjects {
		if projectMatch.Match(prj.Path) {
			allSolutionFiles = append(allSolutionFiles, prj.Files()...)
		} else {
			projectsOutsideSolution = append(projectsOutsideSolution, prj)
		}
	}

	anySolutionFile := NewExactMatch(allSolutionFiles)
	return separateProjects(projectsOutsideSolution, anySolutionFile)
}

func separateProjects(projectsOutsideSolution []*msvc.MsbuildProject, anySolutionFile api.Matcher) ([]string, []string) {
	var lost []string
	var lostWithIncludes []string
	for _, prj := range projectsOutsideSolution {
		if hasFilesIncludedIntoSolution(prj, anySolutionFile) {
			lostWithIncludes = append(lostWithIncludes, prj.Path)
		} else {
			lost = append(lost, prj.Path)
		}
	}
	return lost, lostWithIncludes
}

func hasFilesIncludedIntoSolution(prj *msvc.MsbuildProject, anySolutionFile api.Matcher) bool {
	projectFiles := prj.Files()

	for _, f := range projectFiles {
		if anySolutionFile.Match(f) {
			return true
		}
	}

	return false
}
