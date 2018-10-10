package main

import (
	"fmt"
	"os"
	"path/filepath"
	"solt/solution"
	"strings"
)

type projectSolution struct {
	project  string
	solution string
}

func lostprojectscmd(opt options) error {

	var solutions []string
	folders := readProjectDir(opt.Path, func(we *walkEntry) {
		ext := strings.ToLower(filepath.Ext(we.Name))
		if ext == solutionFileExt {
			sp := filepath.Join(we.Parent, we.Name)
			solutions = append(solutions, sp)
		}
	})

	allProjectsWithinSolutions := getAllSolutionsProjects(solutions)

	projectsOutsideSolution, filesInsideSolution := getOutsideProjectsAndFilesInsideSolution(folders, allProjectsWithinSolutions)

	projectsOutside, projectsOutsideSolutionWithFilesInside := separateOutsideProjects(projectsOutsideSolution, filesInsideSolution)

	sortAndOutputToStdout(projectsOutside)

	if len(projectsOutsideSolutionWithFilesInside) > 0 {
		fmt.Printf("\nThese projects not included into any solution but their files used in projects that included into another projects within a solution.\n")
	}

	sortAndOutputToStdout(projectsOutsideSolutionWithFilesInside)

	var unexistProjects = make(map[string][]string)
	for _, prj := range allProjectsWithinSolutions {
		if _, err := os.Stat(prj.project); os.IsNotExist(err) {
			if found, ok := unexistProjects[prj.solution]; ok {
				found = append(found, prj.project)
				unexistProjects[prj.solution] = found
			} else {
				unexistProjects[prj.solution] = []string{prj.project}
			}
		}
	}

	if len(unexistProjects) > 0 {
		fmt.Printf("\nThese projects included into a solution but not found in the file system.\n")
	}

	outputSortedMapToStdout(unexistProjects, "Solution")

	return nil
}

func getOutsideProjectsAndFilesInsideSolution(folders []*folderInfo, allProjectsWithinSolutions map[string]*projectSolution) ([]*folderInfo, map[string]interface{}) {
	var projectsOutsideSolution []*folderInfo
	var filesInsideSolution = make(map[string]interface{})
	for _, info := range folders {
		if info.project == nil {
			continue
		}

		id := strings.ToUpper(info.project.Id)

		_, ok := allProjectsWithinSolutions[id]
		if !ok {
			projectsOutsideSolution = append(projectsOutsideSolution, info)
		} else {
			filesIncluded := getFilesIncludedIntoProject(info)

			for _, f := range filesIncluded {
				filesInsideSolution[strings.ToUpper(f)] = nil
			}
		}
	}
	return projectsOutsideSolution, filesInsideSolution
}

func separateOutsideProjects(projectsOutsideSolution []*folderInfo, filesInsideSolution map[string]interface{}) ([]string, []string) {
	var projectsOutside []string
	var projectsOutsideSolutionWithFilesInside []string
	for _, info := range projectsOutsideSolution {
		projectFiles := getFilesIncludedIntoProject(info)

		var includedIntoOther = false
		for _, f := range projectFiles {
			pf := strings.ToUpper(f)
			if _, ok := filesInsideSolution[pf]; ok {

				dir := filepath.Dir(*info.projectPath)

				if strings.Contains(pf, strings.ToUpper(dir)) {
					includedIntoOther = true
					break
				}
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
		sln, _ := solution.Parse(solpath)

		for _, p := range sln.Projects {
			// Skip solution folders
			if p.TypeId == "{2150E333-8FDC-42A3-9474-1A3956D46DE8}" {
				continue
			}

			id := strings.ToUpper(p.Id)

			if _, ok := projectsInSolution[id]; !ok {
				parent := filepath.Dir(solpath)
				pp := filepath.Join(parent, p.Path)

				ps := projectSolution{
					project:  pp,
					solution: solpath,
				}

				projectsInSolution[id] = &ps
			}
		}
	}
	return projectsInSolution
}
