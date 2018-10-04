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

func readProjectDir(path string, action func(we *walkEntry)) []*folderInfo {
    ch := make(chan *walkEntry, 1024)

    go func(ch chan<- *walkEntry) {
        walkDirBreadthFirst(path, func(parent string, entry os.FileInfo) {
            if entry.IsDir() {
                return
            }

            ch <- &walkEntry{IsDir: false, Size: entry.Size(), Parent: parent, Name: entry.Name()}
        })
        close(ch)
    }(ch)

    foldersMap := make(map[string]*folderInfo)

    for {
        we, ok := <-ch
        if !ok {
            break
        }

        if we.Name == PackagesConfingFile {
            // Create package model from packages.config
            full := filepath.Join(we.Parent, we.Name)

            pack := Packages{}

            err := unmarshalXml(full, &pack)

            if err != nil {
                log.Print(err)
            }

            info, ok := foldersMap[we.Parent]
            if !ok {
                fi := folderInfo{packages: &pack, projectPath: &full}
                foldersMap[we.Parent] = &fi
            } else {
                info.packages = &pack
            }
        }

        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == CSharpProjectExt || ext == CppProjectExt {

            // Create project model from msbuild project file
            full := filepath.Join(we.Parent, we.Name)
            project := Project{}

            err := unmarshalXml(full, &project)

            if err != nil {
                log.Print(err)
                continue
            }

            info, ok := foldersMap[we.Parent]
            if !ok {
                fi := folderInfo{project: &project, projectPath: &full}
                foldersMap[we.Parent] = &fi
            } else {
                info.project = &project
                info.projectPath = &full
            }
        }

        action(we)
    }

    var result []*folderInfo

    for _, v := range foldersMap {
        result = append(result, v)
    }

    return result
}
