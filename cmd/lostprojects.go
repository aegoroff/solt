package cmd

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"log"
	"os"
	"path/filepath"
	"solt/solution"
	"strings"

	"github.com/spf13/cobra"
)

type projectSolution struct {
	project  string
	solution string
}

// lostprojectsCmd represents the lostprojects command
var lostprojectsCmd = &cobra.Command{
	Use:     "lostprojects",
	Aliases: []string{"p"},
	Short:   "Find projects that not included into any solution",
	Run: func(cmd *cobra.Command, args []string) {
		var solutions []string
		foldersTree := readProjectDir(sourcesPath, func(we *walkEntry) {
			ext := strings.ToLower(filepath.Ext(we.Name))
			if ext == solutionFileExt {
				sp := filepath.Join(we.Parent, we.Name)
				solutions = append(solutions, sp)
			}
		})

		allProjectsWithinSolutions := getAllSolutionsProjects(solutions)

		projectsOutsideSolution, filesInsideSolution := getOutsideProjectsAndFilesInsideSolution(foldersTree, allProjectsWithinSolutions)

		projectsOutside, projectsOutsideSolutionWithFilesInside := separateOutsideProjects(projectsOutsideSolution, filesInsideSolution)

		sortAndOutputToStdout(projectsOutside)

		if len(projectsOutsideSolutionWithFilesInside) > 0 {
			fmt.Printf("\nThese projects are not included into any solution but files from the projects' folders are used in another projects within a solution:\n\n")
		}

		sortAndOutputToStdout(projectsOutsideSolutionWithFilesInside)

		unexistProjects := getUnexistProjects(allProjectsWithinSolutions)

		if len(unexistProjects) > 0 {
			fmt.Printf("\nThese projects are included into a solution but not found in the file system:\n")
		}

		outputSortedMapToStdout(unexistProjects, "Solution")
	},
}

func init() {
	rootCmd.AddCommand(lostprojectsCmd)
}

func getUnexistProjects(allProjectsWithinSolutions map[string]*projectSolution) map[string][]string {
	var result = make(map[string][]string)
	for _, prj := range allProjectsWithinSolutions {
		if _, err := os.Stat(prj.project); !os.IsNotExist(err) {
			continue
		}

		if found, ok := result[prj.solution]; ok {
			found = append(found, prj.project)
			result[prj.solution] = found
		} else {
			result[prj.solution] = []string{prj.project}
		}
	}
	return result
}

func getOutsideProjectsAndFilesInsideSolution(foldersTree *rbtree.RbTree, allProjectsWithinSolutions map[string]*projectSolution) ([]*folderInfo, StringHashSet) {

	var projectsOutsideSolution []*folderInfo
	var filesInsideSolution = make(StringHashSet)

	foldersTree.Ascend(func(c *rbtree.Comparable) bool {
		info := (*c).(projectTreeNode).info
		if info.project == nil {
			return true
		}

		id := strings.ToUpper(info.project.Id)

		_, ok := allProjectsWithinSolutions[id]
		if !ok {
			projectsOutsideSolution = append(projectsOutsideSolution, info)
		} else {
			filesIncluded := getFilesIncludedIntoProject(info)

			for _, f := range filesIncluded {
				filesInsideSolution.Add(strings.ToUpper(f))
			}
		}

		return true
	})

	return projectsOutsideSolution, filesInsideSolution
}

func separateOutsideProjects(projectsOutsideSolution []*folderInfo, filesInsideSolution StringHashSet) ([]string, []string) {
	var projectsOutside []string
	var projectsOutsideSolutionWithFilesInside []string
	for _, info := range projectsOutsideSolution {
		projectFiles := getFilesIncludedIntoProject(info)

		var includedIntoOther = false
		for _, f := range projectFiles {
			pf := strings.ToUpper(f)
			if !filesInsideSolution.Contains(pf) {
				continue
			}

			dir := filepath.Dir(*info.projectPath)

			if strings.Contains(pf, strings.ToUpper(dir)) {
				includedIntoOther = true
				break
			}
		}

		if !includedIntoOther {
			projectsOutside = append(projectsOutside, *info.projectPath)
		} else {
			projectsOutsideSolutionWithFilesInside = append(projectsOutsideSolutionWithFilesInside, *info.projectPath)
		}
	}
	return projectsOutside, projectsOutsideSolutionWithFilesInside
}

func getAllSolutionsProjects(solutions []string) map[string]*projectSolution {

	var projectsInSolution = make(map[string]*projectSolution)
	for _, solpath := range solutions {
		sln, err := solution.Parse(solpath)

		if err != nil {
			log.Println(err)
			continue
		}

		for _, p := range sln.Projects {
			// Skip solution folders
			if p.TypeId == "{2150E333-8FDC-42A3-9474-1A3956D46DE8}" {
				continue
			}

			id := strings.ToUpper(p.Id)

			// Already added
			if _, ok := projectsInSolution[id]; ok {
				continue
			}

			parent := filepath.Dir(solpath)
			pp := filepath.Join(parent, p.Path)

			projectsInSolution[id] = &projectSolution{
				project:  pp,
				solution: solpath,
			}
		}
	}
	return projectsInSolution
}
