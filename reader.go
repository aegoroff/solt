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

type folder struct {
    info *folderInfo
    path string
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

    aggregatech := make(chan *folder, 1024)

    go func(ch <-chan *folder) {
        projectFolders := make(map[string]interface{})
        for {
            f, ok := <-ch
            if !ok {
                break
            }

            if _, ok := projectFolders[f.path]; !ok {
                projectFolders[f.path] = nil
                result = append(result, f.info)
            } else {
                current := result[len(result)-1]

                if current.project == nil {
                    // Project read after packages.config
                    current.project = f.info.project
                    current.projectPath = f.info.projectPath
                } else if current.packages == nil {
                    // Project read before packages.config
                    current.packages = f.info.packages
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
            if f, ok := onPackagesConfig(we); ok {
                aggregatech <- f
            }
        }

        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == CSharpProjectExt || ext == CppProjectExt {
            if f, ok := onMsbuildProject(we); ok {
                aggregatech <- f
            }
        }

        action(we)
    }

    return result
}

func onPackagesConfig(we *walkEntry) (*folder, bool) {
    // Create packages model from packages.config
    full := filepath.Join(we.Parent, we.Name)
    pack := Packages{}
    err := unmarshalXml(full, &pack)
    if err != nil {
        log.Printf("%s: %v\n", full, err)
        return nil, false
    }
    pi := folder{
        info: &folderInfo{packages: &pack, projectPath: &full},
        path: we.Parent,
    }
    return &pi, true
}

func onMsbuildProject(we *walkEntry) (*folder, bool) {
    // Create project model from project file
    full := filepath.Join(we.Parent, we.Name)
    project := Project{}
    err := unmarshalXml(full, &project)
    if err != nil {
        log.Printf("%s: %v\n", full, err)
        return nil, false
    }
    pi := folder{
        info: &folderInfo{project: &project, projectPath: &full},
        path: we.Parent,
    }

    return &pi, true
}
