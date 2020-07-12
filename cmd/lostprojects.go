package cmd

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"path/filepath"
	"solt/internal/sys"
	"solt/msvc"
	"strings"

	"github.com/spf13/cobra"
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

			projectLinksBySolution := make(map[string]collections.StringHashSet, len(solutions))
			// Each found solution
			for _, sln := range solutions {
				links := msvc.SelectAllSolutionProjectPaths(sln, func(s string) string { return s })
				projectLinksBySolution[sln.Path] = links
				// to create projectsInSolutions you shoud normalize path to build Matcher
				for _, item := range links.ItemsDecorated(normalize) {
					linkedProjects = append(linkedProjects, item)
				}
			}

			lost, lostWithIncludes := findLostProjects(allProjects, linkedProjects)

			// Lost projects
			sortAndOutput(appWriter, lost)

			if len(lostWithIncludes) > 0 {
				m := "\n<red>These projects are not included into any solution but files from the projects' folders are used in another projects within a solution:</>\n\n"
				color.Fprintf(appWriter, m)
			}

			// Lost projects that have includes files that used
			sortAndOutput(appWriter, lostWithIncludes)

			unexistProjects := getUnexistProjects(projectLinksBySolution, appFileSystem)

			if len(unexistProjects) > 0 {
				color.Fprintf(appWriter, "\n<red>These projects are included into a solution but not found in the file system:</>\n")
			}

			// Included but not exist in FS
			outputSortedMap(appWriter, unexistProjects, "Solution")

			return nil
		},
	}

	return cmd
}

func getUnexistProjects(projectsInSolutions map[string]collections.StringHashSet, fs afero.Fs) map[string][]string {
	var result = make(map[string][]string, len(projectsInSolutions))

	filer := sys.NewFiler(fs, appWriter)
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
	var filesInsideSolution = make(collections.StringHashSet, len(allProjects)*20)

	for _, prj := range allProjects {
		// Path in upper registry is the project's key
		projectKey := normalize(prj.Path)

		ok := projectMatch.Match(projectKey)
		if !ok {
			projectsOutsideSolution = append(projectsOutsideSolution, prj)
		} else {
			filesIncluded := msvc.GetFilesIncludedIntoProject(prj)

			for _, f := range filesIncluded {
				filesInsideSolution.Add(normalize(f))
			}
		}
	}

	return separateProjects(projectsOutsideSolution, filesInsideSolution)
}

func separateProjects(projectsOutsideSolution []*msvc.MsbuildProject, filesInsideSolution collections.StringHashSet) ([]string, []string) {
	var lost []string
	var lostWithIncludes []string
	for _, prj := range projectsOutsideSolution {
		includedIntoOther := hasFilesIncludedIntoActual(prj, filesInsideSolution)

		if !includedIntoOther {
			lost = append(lost, prj.Path)
		} else {
			lostWithIncludes = append(lostWithIncludes, prj.Path)
		}
	}
	return lost, lostWithIncludes
}

func hasFilesIncludedIntoActual(prj *msvc.MsbuildProject, solutionFiles collections.StringHashSet) bool {
	projectFiles := msvc.GetFilesIncludedIntoProject(prj)

	pdir := filepath.Dir(prj.Path)

	for _, f := range projectFiles {
		pfile := normalize(f)
		if !solutionFiles.Contains(pfile) {
			continue
		}

		if strings.HasPrefix(pfile, normalize(pdir)) {
			return true
		}
	}
	return false
}
