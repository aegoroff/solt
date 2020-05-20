package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"solt/solution"
	"strings"

	"github.com/spf13/cobra"
)

type projectSolution struct {
	id       string
	path     string
	solution string
}

// lostprojectsCmd represents the lostprojects command
var lostprojectsCmd = &cobra.Command{
	Use:     "lostprojects",
	Aliases: []string{"lp"},
	Short:   "Find projects that not included into any solution",
	Run: func(cmd *cobra.Command, args []string) {

		foldersTree := readProjectDir(sourcesPath, appFileSystem, func(we *walkEntry) {})

		allProjectsWithinSolutions := getAllSolutionsProjects(foldersTree)

		projectsOutsideSolution, filesInsideSolution := getOutsideProjectsAndFilesInsideSolution(foldersTree, allProjectsWithinSolutions)

		projectsOutside, projectsOutsideSolutionWithFilesInside := separateOutsideProjects(projectsOutsideSolution, filesInsideSolution)

		sortAndOutput(appWriter, projectsOutside)

		if len(projectsOutsideSolutionWithFilesInside) > 0 {
			fmt.Fprintf(appWriter, "\nThese projects are not included into any solution but files from the projects' folders are used in another projects within a solution:\n\n")
		}

		sortAndOutput(appWriter, projectsOutsideSolutionWithFilesInside)

		unexistProjects := getUnexistProjects(allProjectsWithinSolutions, appFileSystem)

		if len(unexistProjects) > 0 {
			fmt.Fprintf(appWriter, "\nThese projects are included into a solution but not found in the file system:\n")
		}

		outputSortedMap(appWriter, unexistProjects, "Solution")
	},
}

func init() {
	rootCmd.AddCommand(lostprojectsCmd)
}

func getUnexistProjects(allProjectsWithinSolutions map[string]*projectSolution, fs afero.Fs) map[string][]string {
	var result = make(map[string][]string)
	for _, prj := range allProjectsWithinSolutions {
		if _, err := fs.Stat(prj.path); !os.IsNotExist(err) {
			continue
		}

		if found, ok := result[prj.solution]; ok {
			found = append(found, prj.path)
			result[prj.solution] = found
		} else {
			result[prj.solution] = []string{prj.path}
		}
	}
	return result
}

func getOutsideProjectsAndFilesInsideSolution(foldersTree *rbtree.RbTree, allProjectsWithinSolutions map[string]*projectSolution) ([]*msbuildProject, collections.StringHashSet) {
	var projectsOutsideSolution []*msbuildProject
	var filesInsideSolution = make(collections.StringHashSet)

	foldersTree.Ascend(func(c *rbtree.Comparable) bool {
		folder := (*c).(*folder)
		content := folder.content
		if len(content.projects) == 0 {
			return true
		}

		for _, prj := range content.projects {
			// Path in upper registry is the project's key
			projectKey := strings.ToUpper(prj.path)

			_, ok := allProjectsWithinSolutions[projectKey]
			if !ok {
				projectsOutsideSolution = append(projectsOutsideSolution, prj)
			} else {
				filesIncluded := getFilesIncludedIntoProject(prj)

				for _, f := range filesIncluded {
					filesInsideSolution.Add(strings.ToUpper(f))
				}
			}
		}

		return true
	})

	return projectsOutsideSolution, filesInsideSolution
}

func separateOutsideProjects(projectsOutsideSolution []*msbuildProject, filesInsideSolution collections.StringHashSet) ([]string, []string) {
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

func getAllSolutionsProjects(foldersTree *rbtree.RbTree) map[string]*projectSolution {

	var projectsInSolution = make(map[string]*projectSolution)

	// Select only folders that contain solution(s)
	foldersTree.WalkInorder(func(n *rbtree.Node) {
		fold := (*n.Key).(*folder)
		content := fold.content
		if len(content.solutions) == 0 {
			return
		}

		for _, sln := range content.solutions {

			for _, prj := range sln.solution.Projects {
				// Skip solution folders
				if prj.TypeId == solution.IdSolutionFolder {
					continue
				}

				fullProjectPath := filepath.Join(fold.path, prj.Path)
				key := strings.ToUpper(fullProjectPath)

				// Already added
				if _, ok := projectsInSolution[key]; ok {
					continue
				}

				projectsInSolution[key] = &projectSolution{
					path:     fullProjectPath,
					id:       strings.ToUpper(prj.Id),
					solution: sln.path,
				}
			}
		}
	})

	return projectsInSolution
}
