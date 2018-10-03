package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "text/tabwriter"
)

type mismatch struct {
    pkg      string
    versions []string
}

func nugetcmd(opt options) error {
    var solutions []string
    foldersMap := readProjectDir(opt.Path, func(we *walkEntry) {
        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == SolutionFileExt {
            solutions = append(solutions, filepath.Join(we.Parent, we.Name))
        }
    })

    if opt.Nuget.Mismatch {
        showMismatches(solutions, foldersMap)
    } else {
        showPackagesInfoByFolders(foldersMap)
    }

    return nil
}

func showMismatches(solutions []string, foldersMap map[string]*folderInfo) {

    solutionProjects := getProjectsOfSolutions(solutions, foldersMap)

    mismatches := calculateMismatches(solutionProjects)

    if len(mismatches) == 0 {
        return
    }

    fmt.Println(" Found different nuget package's versions in the same solution:")

    const format = "  %v\t%v\n"
    tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

    fmt.Fprintf(tw, format, "Package", "Versions")
    fmt.Fprintf(tw, format, "-------", "--------")

    for sol, m := range mismatches {
        fmt.Printf("\n %s\n", sol)
        for _, item := range m {
            fmt.Fprintf(tw, format, item.pkg, strings.Join(item.versions, ", "))
        }
        tw.Flush()
    }
}

func getProjectsOfSolutions(solutions []string, foldersMap map[string]*folderInfo) map[string][]*folderInfo {
    var solutionProjects = make(map[string][]*folderInfo)
    for _, sol := range solutions {
        projects, _, _ := parseSolution(sol)
        var solutionProjectIds = make(map[string]interface{})
        for _, sp := range projects {
            solutionProjectIds[sp.Id] = nil
        }

        for _, finfo := range foldersMap {
            if finfo.project == nil {
                continue
            }

            if _, ok := solutionProjectIds[finfo.project.Id]; ok {
                if v, ok := solutionProjects[sol]; !ok {
                    solutionProjects[sol] = []*folderInfo{finfo}
                } else {
                    solutionProjects[sol] = append(v, finfo)
                }
            }
        }
    }
    return solutionProjects
}

func calculateMismatches(solutionProjects map[string][]*folderInfo) map[string][]*mismatch {
    var mismatches = make(map[string][]*mismatch)
    for solution, projects := range solutionProjects {
        var packagesMap = make(map[string][]string)
        for _, prj := range projects {
            if prj.packages == nil {
                continue
            }

            for _, pkg := range prj.packages.Packages {
                if v, ok := packagesMap[pkg.Id]; !ok {
                    packagesMap[pkg.Id] = []string{pkg.Version}
                } else {
                    if contains(v, pkg.Version) {
                        continue
                    }

                    packagesMap[pkg.Id] = append(v, pkg.Version)
                }
            }
        }

        for pkg, vers := range packagesMap {
            if len(vers) < 2 {
                continue
            }

            m := mismatch{
                pkg:      pkg,
                versions: vers,
            }

            if v, ok := mismatches[solution]; !ok {
                mismatches[solution] = []*mismatch{&m}
            } else {
                mismatches[solution] = append(v, &m)
            }
        }
    }
    return mismatches
}

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func showPackagesInfoByFolders(foldersMap map[string]*folderInfo) {
    const format = "  %v\t%v\n"
    tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 4, ' ', 0)

    for k, v := range foldersMap {
        if v.packages == nil {
            continue
        }
        fmt.Printf(" %s\n", k)
        fmt.Fprintf(tw, format, "Package", "Version")
        fmt.Fprintf(tw, format, "-------", "--------")

        for _, p := range v.packages.Packages {
            fmt.Fprintf(tw, format, p.Id, p.Version)
        }
        tw.Flush()
        fmt.Println()
    }
}
