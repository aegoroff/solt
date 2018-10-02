package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "solt/solution"
    "strings"
)

func analyzecmd(opt options) error {

    solutionPath := filepath.Join(opt.Path, opt.Analyze.Solution)

    collectSolutions := opt.Analyze.Solution == ""

    var solutions []string
    foldersMap := readProjectDir(opt.Path, func(we *walkEntry) {
        if !collectSolutions {
            return
        }
        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == SolutionFileExt {
            solutions = append(solutions, filepath.Join(we.Parent, we.Name))
        }
    })

    if collectSolutions {
        var projectsInSolution = make(map[string]*solution.Project)
        for _, sol := range solutions {
            projects, _, _ := parseSolution(sol)
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
    } else {
        return analyzeSolution(solutionPath, foldersMap)
    }
}

func analyzeSolution(solutionPath string, foldersMap map[string]*folderInfo) error {

    projects, globalSections, err := parseSolution(solutionPath)

    if err != nil {
        return err
    }

    var projectsInSolution = make(map[string]*solution.Project)
    for _, p := range projects {
        projectsInSolution[p.Id] = p
    }

    fmt.Printf("Projects: %d Global Sections: %d\n", len(projects), len(globalSections))

    for _, p := range foldersMap {
        if p.project == nil {
            continue
        }
        _, ok := projectsInSolution[p.project.Id]
        if !ok {
            fmt.Printf("Project %s not included into solution [%s]\n", *p.projectPath, solutionPath)
        }

    }

    return nil
}

func parseSolution(solutionPath string) ([]*solution.Project, []*solution.Section, error) {

    f, err := os.Open(solutionPath)
    if err != nil {
        return nil, nil, err
    }
    defer f.Close()

    br := bufio.NewReader(f)
    r, _, err := br.ReadRune()
    if err != nil {
        return nil, nil, err
    }
    if r != '\uFEFF' {
        br.UnreadRune() // Not a BOM -- put the rune back
    }

    bs := bufio.NewScanner(br)
    bs.Split(bufio.ScanRunes)
    sb := strings.Builder{}

    for bs.Scan() {
        sb.WriteString(bs.Text())
    }

    str := sb.String()

    projects, globalSections := solution.Parse(str)

    return projects, globalSections, nil
}
