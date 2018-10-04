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
            solutions = append(solutions, filepath.Join(we.Parent, we.Name))
        }
    })

    var projectsInSolution = make(map[string]interface{})
    for _, sol := range solutions {
        sln, _ := solution.Parse(sol)
        parent := filepath.Dir(sol)

        for _, p := range sln.Projects {
            pp := filepath.Join(parent, p.Path)

            if _, ok := projectsInSolution[pp]; !ok {
                projectsInSolution[pp] = nil
            }
        }
    }

    for _, p := range foldersMap {
        if p.project == nil {
            continue
        }
        _, ok := projectsInSolution[*p.projectPath]
        if !ok {
            fmt.Println(*p.projectPath)
        }
    }

    return nil
}
