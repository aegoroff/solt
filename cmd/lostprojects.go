package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"solt/internal/msvc"
	"strings"

	"github.com/spf13/cobra"
)

// lostprojectsCmd represents the lostprojects command
var lostprojectsCmd = &cobra.Command{
	Use:     "lp",
	Aliases: []string{"lostprojects"},
	Short:   "Find projects that not included into any solution",
	RunE: func(cmd *cobra.Command, args []string) error {

		foldersTree := msvc.ReadSolutionDir(sourcesPath, appFileSystem)

		solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

		// linked from any solution projects list
		// so these projects are not considered lost
		var linkedProjects []string

		projectLinksBySolution := make(map[string]collections.StringHashSet)
		// Each found solution
		for _, sln := range solutions {
			links := msvc.SelectAllSolutionProjectPaths(sln, func(s string) string { return s })
			projectLinksBySolution[sln.Path] = links
			// to create projectsInSolutions you shoud normalize path to build Matcher
			for _, item := range links.ItemsDecorated(normalize) {
				linkedProjects = append(linkedProjects, item)
			}
		}

		projectsOutsideSolutions, filesInsideSolution := filterProjects(allProjects, linkedProjects)
		lostProjects, lostProjectsThatIncludeSolutionProjectsFiles := separateProjects(projectsOutsideSolutions, filesInsideSolution)

		sortAndOutput(appWriter, lostProjects)

		if len(lostProjectsThatIncludeSolutionProjectsFiles) > 0 {
			_, _ = fmt.Fprintf(appWriter, "\nThese projects are not included into any solution but files from the projects' folders are used in another projects within a solution:\n\n")
		}

		sortAndOutput(appWriter, lostProjectsThatIncludeSolutionProjectsFiles)

		unexistProjects := getUnexistProjects(projectLinksBySolution, appFileSystem)

		if len(unexistProjects) > 0 {
			_, _ = fmt.Fprintf(appWriter, "\nThese projects are included into a solution but not found in the file system:\n")
		}

		outputSortedMap(appWriter, unexistProjects, "Solution")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(lostprojectsCmd)
}

func getUnexistProjects(projectsInSolutions map[string]collections.StringHashSet, fs afero.Fs) map[string][]string {
	var result = make(map[string][]string)
	for sol, projects := range projectsInSolutions {
		for _, prj := range projects.Items() {
			if _, err := fs.Stat(prj); !os.IsNotExist(err) {
				continue
			}

			if found, ok := result[sol]; ok {
				found = append(found, prj)
				result[sol] = found
			} else {
				result[sol] = []string{prj}
			}
		}
	}
	return result
}

func filterProjects(allProjects []*msvc.MsbuildProject, linkedProjects []string) ([]*msvc.MsbuildProject, collections.StringHashSet) {
	// Create projects matching machine
	projectMatch := NewExactMatchS(linkedProjects)
	var projectsOutsideSolution []*msvc.MsbuildProject
	var filesInsideSolution = make(collections.StringHashSet)

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

	return projectsOutsideSolution, filesInsideSolution
}

func separateProjects(projectsOutsideSolution []*msvc.MsbuildProject, filesInsideSolution collections.StringHashSet) ([]string, []string) {
	var projectsOutside []string
	var projectsOutsideSolutionWithFilesInside []string
	for _, prj := range projectsOutsideSolution {
		projectFiles := msvc.GetFilesIncludedIntoProject(prj)

		var includedIntoOther = false
		for _, f := range projectFiles {
			pf := normalize(f)
			if !filesInsideSolution.Contains(pf) {
				continue
			}

			dir := filepath.Dir(prj.Path)

			if strings.Contains(pf, normalize(dir)) {
				includedIntoOther = true
				break
			}
		}

		if !includedIntoOther {
			projectsOutside = append(projectsOutside, prj.Path)
		} else {
			projectsOutsideSolutionWithFilesInside = append(projectsOutsideSolutionWithFilesInside, prj.Path)
		}
	}
	return projectsOutside, projectsOutsideSolutionWithFilesInside
}
