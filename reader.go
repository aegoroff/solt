package main

import (
    "log"
    "os"
    "path/filepath"
    "strings"
)

type folderInfo struct {
    packages    *Packages
    project     *Project
    projectPath *string
}

func getFilesIncludedIntoProject(info *folderInfo) []string {
    dir := filepath.Dir(*info.projectPath)
    var result []string
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

type folderInfoPair struct {
    info   *folderInfo
    parent string
}

func readProjectDir(path string, action func(we *walkEntry)) []*folderInfo {
    readch := make(chan *walkEntry, 1024)

    go func(ch chan<- *walkEntry) {
        walkDirBreadthFirst(path, func(parent string, entry os.FileInfo) {
            if entry.IsDir() {
                return
            }

            ch <- &walkEntry{IsDir: false, Size: entry.Size(), Parent: parent, Name: entry.Name()}
        })
        close(ch)
    }(readch)

    var result []*folderInfo
    parents := make(map[string]interface{})

    aggregatech := make(chan *folderInfoPair, 1024)

    go func(ch <-chan *folderInfoPair) {
        for {
            fi, ok := <-ch
            if !ok {
                break
            }

            if _, ok := parents[fi.parent]; !ok {
                parents[fi.parent] = nil
                result = append(result, fi.info)
            } else {
                current := result[len(result)-1]

                if current.project == nil {
                    // Project read after packages.config
                    current.project = fi.info.project
                    current.projectPath = fi.info.projectPath
                } else if current.packages == nil {
                    // Project read before packages.config
                    current.packages = fi.info.packages
                }
            }
        }
    }(aggregatech)

    for {
        we, ok := <-readch
        if !ok {
            close(aggregatech)
            break
        }

        if we.Name == PackagesConfigFile {
            // Create package model from packages.config
            full := filepath.Join(we.Parent, we.Name)
            pack := Packages{}
            err := unmarshalXml(full, &pack)

            if err != nil {
                log.Printf("%s: %v\n", full, err)
                continue
            }

            pi := folderInfoPair{
                info:   &folderInfo{packages: &pack, projectPath: &full},
                parent: we.Parent,
            }

            aggregatech <- &pi
        }

        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == CSharpProjectExt || ext == CppProjectExt {
            // Create project model from msbuild project file
            full := filepath.Join(we.Parent, we.Name)
            project := Project{}
            err := unmarshalXml(full, &project)

            if err != nil {
                log.Printf("%s: %v\n", full, err)
                continue
            }

            pi := folderInfoPair{
                info:   &folderInfo{project: &project, projectPath: &full},
                parent: we.Parent,
            }

            aggregatech <- &pi
        }

        action(we)
    }

    return result
}
