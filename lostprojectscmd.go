package main

import (
    "bufio"
    "fmt"
    "os"
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
