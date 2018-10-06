package main

import (
	"fmt"
	"path/filepath"
	"solt/solution"
	"sort"
	"strings"
)

func lostprojectscmd(opt options) error {

	var solutions []string
	folders := readProjectDir(opt.Path, func(we *walkEntry) {
		ext := strings.ToLower(filepath.Ext(we.Name))
		if ext == solutionFileExt {
			sp := filepath.Join(we.Parent, we.Name)
			solutions = append(solutions, sp)
		}
	})

	var projectsInSolution = make(map[string]interface{})
	for _, solpath := range solutions {
		sln, _ := solution.Parse(solpath)

		parent := filepath.Dir(solpath)

		for _, p := range sln.Projects {
			// Skip solution folders
			if p.TypeId == "{2150E333-8FDC-42A3-9474-1A3956D46DE8}" {
				continue
			}

			pp := filepath.Join(parent, p.Path)

			if _, ok := projectsInSolution[pp]; !ok {
				projectsInSolution[pp] = nil
			}
		}
	}

	var projectsOutsideSolution []*folderInfo
	var filesInsideSolution = make(map[string]interface{})
	for _, info := range folders {
		if info.project == nil {
			continue
		}
		project := *info.projectPath
		_, ok := projectsInSolution[project]
		if !ok {
			projectsOutsideSolution = append(projectsOutsideSolution, info)
		} else {
			filesIncluded := getFilesIncludedIntoProject(info)

			for _, f := range filesIncluded {
				filesInsideSolution[f] = nil
			}
		}
	}

	var projectsOutsideSolutionWithFilesInside []string
	for _, info := range projectsOutsideSolution {
		filesIncluded := getFilesIncludedIntoProject(info)

		var includedIntoOther = false
		for _, f := range filesIncluded {
			if _, ok := filesInsideSolution[f]; ok {

				dir := filepath.Dir(*info.projectPath)

				if strings.Contains(f, dir) {
					includedIntoOther = true
					break
				}
			}
		}

		if !includedIntoOther {
			fmt.Printf(" %s\n", *info.projectPath)
		} else {
			projectsOutsideSolutionWithFilesInside = append(projectsOutsideSolutionWithFilesInside, *info.projectPath)
		}
	}

	if len(projectsOutsideSolutionWithFilesInside) > 0 {
		fmt.Printf("\nThese projects not included into any solution but their files used in projects that included into another projects within a solution.\n")
	}

	sort.Strings(projectsOutsideSolutionWithFilesInside)

	for _, p := range projectsOutsideSolutionWithFilesInside {
		fmt.Printf(" %s\n", p)
	}

	return nil
}
