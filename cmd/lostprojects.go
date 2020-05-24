package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	goahocorasick "github.com/anknown/ahocorasick"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// lostprojectsCmd represents the lostprojects command
var lostprojectsCmd = &cobra.Command{
	Use:     "lp",
	Aliases: []string{"lostprojects"},
	Short:   "Find projects that not included into any solution",
	RunE: func(cmd *cobra.Command, args []string) error {

		foldersTree := readProjectDir(sourcesPath, appFileSystem, func(we *walkEntry) {})

		solutions := selectSolutions(foldersTree)

		var projectsInSolutions []string
		projectsBySolution := make(map[string]collections.StringHashSet)
		// Each found solution
		for _, sln := range solutions {
			solutionProjectPaths := selectAllSolutionProjectPaths(sln, false)
			projectsBySolution[sln.path] = solutionProjectPaths
			for _, item := range solutionProjectPaths.Items() {
				projectsInSolutions = append(projectsInSolutions, strings.ToUpper(item))
			}
		}

		// Create projects matching machine
		pmm, err := createAhoCorasickMachine(projectsInSolutions)
		if err != nil {
			return err
		}

		projectsOutsideSolutions, filesInsideSolution := getOutsideProjectsAndFilesInsideSolution(foldersTree, pmm)

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

		if showMemUsage {
			printMemUsage(appWriter)
		}

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

func getOutsideProjectsAndFilesInsideSolution(foldersTree *rbtree.RbTree, pmm *goahocorasick.Machine) ([]*msbuildProject, collections.StringHashSet) {
	var projectsOutsideSolution []*msbuildProject
	var filesInsideSolution = make(collections.StringHashSet)

	walkProjects(foldersTree, func(prj *msbuildProject, fold *folder) {
		// Path in upper registry is the project's key
		projectKey := strings.ToUpper(prj.path)

		ok := Match(pmm, projectKey)
		if !ok {
			projectsOutsideSolution = append(projectsOutsideSolution, prj)
		} else {
			filesIncluded := getFilesIncludedIntoProject(prj)

			for _, f := range filesIncluded {
				filesInsideSolution.Add(strings.ToUpper(f))
			}
		}
	})

	return projectsOutsideSolution, filesInsideSolution
}

func separateProjects(projectsOutsideSolution []*msbuildProject, filesInsideSolution collections.StringHashSet) ([]string, []string) {
	var projectsOutside []string
	var projectsOutsideSolutionWithFilesInside []string
	for _, prj := range projectsOutsideSolution {
		projectFiles := getFilesIncludedIntoProject(prj)

		var includedIntoOther = false
		for _, f := range projectFiles {
			pf := strings.ToUpper(f)
			if !filesInsideSolution.Contains(pf) {
				continue
			}

			dir := filepath.Dir(prj.path)

			if strings.Contains(pf, strings.ToUpper(dir)) {
				includedIntoOther = true
				break
			}
		}

		if !includedIntoOther {
			projectsOutside = append(projectsOutside, prj.path)
		} else {
			projectsOutsideSolutionWithFilesInside = append(projectsOutsideSolutionWithFilesInside, prj.path)
		}
	}
	return projectsOutside, projectsOutsideSolutionWithFilesInside
}
