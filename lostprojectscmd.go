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

    var projectsInSolution = make(map[string]*solution.Project)
    for _, sol := range solutions {
        projects, _, _ := solution.Parse(sol)
        for _, p := range projects {
            if _, ok := projectsInSolution[p.Id]; !ok {
                projectsInSolution[p.Id] = p
            }
        }
    }

    for _, p := range foldersMap {
        if p.project == nil {
            continue
        }
        _, ok := projectsInSolution[p.project.Id]
        if !ok {
            fmt.Println(*p.projectPath)
        }
    }

    return nil
}
