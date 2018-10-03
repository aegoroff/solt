package main

import (
    "fmt"
    "path/filepath"
    "strings"
)

var subfolderToExclude = []string{
    "obj",
    "bin",
}

func lostfilescmd(opt options) error {

    filter := CSharpCodeFileExt
    if len(opt.LostFiles.Filter) > 0 {
        filter = opt.LostFiles.Filter
    }

    var foundFiles []string
    var packagesFolders = make(map[string]interface{})
    foldersMap := readProjectDir(opt.Path, func(we *walkEntry) {
        // Add file to filtered files slice
        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == filter {
            fp := filepath.Join(we.Parent, we.Name)
            foundFiles = append(foundFiles, fp)
        }

        if ext == SolutionFileExt {
            ppath := filepath.Join(we.Parent, "packages")
            if _, ok := packagesFolders[ppath]; !ok {
                packagesFolders[ppath] = nil
            }
        }
    })

    includedFiles, excludedFolders := createIncludedFilesAndExcludedFolders(foldersMap)

    for k := range packagesFolders {
        excludedFolders = append(excludedFolders, k)
    }

    findLostFiles(excludedFolders, foundFiles, includedFiles)

    return nil
}

func findLostFiles(excludedFolders []string, foundFiles []string, includedFiles map[string]interface{}) {
    excludedFoldersMachine := createAhoCorasickMachine(excludedFolders)
    for _, file := range foundFiles {
        if _, ok := includedFiles[file]; !ok {

            if Match(excludedFoldersMachine, file) {
                continue
            }

            fmt.Println(file)
        }
    }
}

func createIncludedFilesAndExcludedFolders(foldersMap map[string]*folderInfo) (map[string]interface{}, []string) {
    var excludeFolders []string
    var includedFiles = make(map[string]interface{})
    for parent, info := range foldersMap {
        if info.project == nil {
            continue
        }

        // Add project base + exclude subfolder into exclude folders list
        for _, s := range subfolderToExclude {
            sub := filepath.Join(parent, s)
            excludeFolders = append(excludeFolders, sub)
        }

        // Add compiles, contents and nones into included files map
        addFilesFrom(info.project.Compiles, parent, includedFiles)
        addFilesFrom(info.project.Contents, parent, includedFiles)
        addFilesFrom(info.project.Nones, parent, includedFiles)
        addFilesFrom(info.project.CLCompiles, parent, includedFiles)
        addFilesFrom(info.project.CLInclude, parent, includedFiles)
    }
    return includedFiles, excludeFolders
}

func addFilesFrom(src []Include, parent string, includedFiles map[string]interface{}) {
    if src == nil {
        return
    }

    for _, c := range src {
        fp := filepath.Join(parent, c.Path)
        includedFiles[fp] = nil
    }
}
