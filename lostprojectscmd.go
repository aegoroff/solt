package main

import (
    "fmt"
    "path/filepath"
    "solt/solution"
    "strings"
)

func lostprojectscmd(opt options) error {

    var solutions []string
    foldersMap := readProjectDir(opt.Path, func(we *walkEntry) {
        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == SolutionFileExt {
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
    for _, info := range foldersMap {
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
                includedIntoOther = true
                break
            }
        }

        if !includedIntoOther {
            fmt.Printf(" %s\n", *info.projectPath)
        } else {
            projectsOutsideSolutionWithFilesInside = append(projectsOutsideSolutionWithFilesInside, *info.projectPath)
        }
    }

    if len(projectsOutsideSolutionWithFilesInside) > 0 {
        fmt.Printf("\nThese projects not included into any solution but their files used in projects that included into another projects within solution.\n")
    }

    for _, p := range projectsOutsideSolutionWithFilesInside {
        fmt.Printf(" %s\n", p)
    }

    return nil
}

func getFilesIncludedIntoProject(info *folderInfo) []string {
    dir := filepath.Dir(*info.projectPath)
    var  result []string
    result = append(result, getFiles(info.project.Contents, dir)...)
    result = append(result, getFiles(info.project.Nones, dir)...)
    result = append(result, getFiles(info.project.CLCompiles, dir)...)
    result = append(result, getFiles(info.project.CLInclude, dir)...)
    result = append(result, getFiles(info.project.Compiles, dir)...)
    return result
}

func getFiles(includes []Include, dir string) []string {
    if includes == nil {
        return []string{}
    }

    var result []string

    for _, c := range includes {
        fp := filepath.Join(dir, c.Path)
        result = append(result, fp)
    }

    return result
}
