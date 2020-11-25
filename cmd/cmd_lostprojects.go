package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"solt/internal/sys"
	"solt/msvc"
)

func newLostProjects() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "lp",
		Aliases: []string{"lostprojects"},
		Short:   "Find projects that not included into any solution",
		RunE: func(cmd *cobra.Command, args []string) error {
			foldersTree := msvc.ReadSolutionDir(sourcesPath, appFileSystem)

			solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

			// linked from any solution projects list
			// so these projects are not considered lost
			var linkedProjects []string

			projectLinksBySolution := make(map[string]c9s.StringHashSet)
			// Each found solution
			for _, sln := range solutions {
				links := msvc.SelectAllSolutionProjectPaths(sln, func(s string) string { return s })
				projectLinksBySolution[sln.Path] = links
				// to create projectsInSolutions you should normalize path to build Matcher
				for _, item := range links.ItemsDecorated(normalize) {
					linkedProjects = append(linkedProjects, item)
				}
			}

			lost, lostWithIncludes := findLostProjects(allProjects, linkedProjects)

			// Lost projects
			sortAndOutput(appPrinter, lost)

			if len(lostWithIncludes) > 0 {
				m := "\n<red>These projects are not included into any solution but files from the projects' folders are used in another projects within a solution:</>\n\n"
				appPrinter.cprint(m)
			}

			// Lost projects that have includes files that used
			sortAndOutput(appPrinter, lostWithIncludes)

			unexistProjects := getUnexistProjects(projectLinksBySolution, appFileSystem)

			if len(unexistProjects) > 0 {
				appPrinter.cprint("\n<red>These projects are included into a solution but not found in the file system:</>\n")
			}

			// Included but not exist in FS
			outputSortedMap(appPrinter, unexistProjects, "Solution")

			return nil
		},
	}

	return cmd
}

func getUnexistProjects(projectsInSolutions map[string]c9s.StringHashSet, fs afero.Fs) map[string][]string {
	var result = make(map[string][]string)

	filer := sys.NewFiler(fs, appPrinter.writer())
	for spath, projects := range projectsInSolutions {
		nonexist := filer.CheckExistence(projects.Items())

		if len(nonexist) > 0 {
			result[spath] = append(result[spath], nonexist...)
		}
	}
	return result
}

func findLostProjects(allProjects []*msvc.MsbuildProject, linkedProjects []string) ([]string, []string) {
	// Create projects matching machine
	projectMatch := NewExactMatchS(linkedProjects)
	var projectsOutsideSolution []*msvc.MsbuildProject
	var allSolutionFiles = make(c9s.StringHashSet)

	for _, prj := range allProjects {
		// Path in upper registry is the project's key
		projectKey := normalize(prj.Path)

		ok := projectMatch.Match(projectKey)
		if !ok {
			projectsOutsideSolution = append(projectsOutsideSolution, prj)
		} else {
			filesIncluded := msvc.GetFilesIncludedIntoProject(prj)

			for _, f := range filesIncluded {
				allSolutionFiles.Add(normalize(f))
			}
		}
	}

	return separateProjects(projectsOutsideSolution, allSolutionFiles)
}

func separateProjects(projectsOutsideSolution []*msvc.MsbuildProject, allSolutionFiles c9s.StringHashSet) ([]string, []string) {
	var lost []string
	var lostWithIncludes []string
	for _, prj := range projectsOutsideSolution {
		if hasFilesIncludedIntoSolution(prj, allSolutionFiles) {
			lostWithIncludes = append(lostWithIncludes, prj.Path)
		} else {
			lost = append(lost, prj.Path)
		}
	}
	return lost, lostWithIncludes
}

func hasFilesIncludedIntoSolution(prj *msvc.MsbuildProject, allSolutionFiles c9s.StringHashSet) bool {
	projectFiles := msvc.GetFilesIncludedIntoProject(prj)

	for _, f := range projectFiles {
		pfile := normalize(f)
		if allSolutionFiles.Contains(pfile) {
			return true
		}
	}

	return false
}
