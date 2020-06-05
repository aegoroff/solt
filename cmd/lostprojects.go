package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
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

		solutions := msvc.SelectSolutions(foldersTree)

		var projectsInSolutions []string
		projectsBySolution := make(map[string]collections.StringHashSet)
		// Each found solution
		for _, sln := range solutions {
			solutionProjectPaths := msvc.SelectAllSolutionProjectPaths(sln, func(s string) string { return s })
			projectsBySolution[sln.Path] = solutionProjectPaths
			// to create projectsInSolutions you shoud normalize path to build AhoCorasickMachine
			for _, item := range solutionProjectPaths.ItemsDecorated(normalize) {
				projectsInSolutions = append(projectsInSolutions, item)
			}
		}

		// Create projects matching machine
		matcher, err := NewPartialMatcher(projectsInSolutions)
		if err != nil {
			return err
		}

		projectsOutsideSolutions, filesInsideSolution := getOutsideProjectsAndFilesInsideSolution(foldersTree, matcher)

		lostProjects, lostProjectsThatIncludeSolutionProjectsFiles := separateProjects(projectsOutsideSolutions, filesInsideSolution)

		sortAndOutput(appWriter, lostProjects)

		if len(lostProjectsThatIncludeSolutionProjectsFiles) > 0 {
			_, _ = fmt.Fprintf(appWriter, "\nThese projects are not included into any solution but files from the projects' folders are used in another projects within a solution:\n\n")
		}

		sortAndOutput(appWriter, lostProjectsThatIncludeSolutionProjectsFiles)

		unexistProjects := getUnexistProjects(projectsBySolution, appFileSystem)

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

func getOutsideProjectsAndFilesInsideSolution(ftree rbtree.RbTree, projectMatch Matcher) ([]*msvc.MsbuildProject, collections.StringHashSet) {
	var projectsOutsideSolution []*msvc.MsbuildProject
	var filesInsideSolution = make(collections.StringHashSet)

	projects := msvc.SelectProjects(ftree)
	for _, prj := range projects {
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
