package main

import (
    "fmt"
    "github.com/anknown/ahocorasick"
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
    foldersMap := readProjectDir(opt.Path, func(we *walkEntry) {
        // Add file to filtered files slice
        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == filter {
            fp := filepath.Join(we.Parent, we.Name)
            foundFiles = append(foundFiles, fp)
        }
    })

    includedFiles, excludeFoldersMachine := createFilesAndFoldersMap(foldersMap)

    for _, file := range foundFiles {
        if _, ok := includedFiles[file]; !ok {

            if Match(excludeFoldersMachine, file) {
                continue
            }

            fmt.Println(file)
        }
    }

    return nil
}

func createFilesAndFoldersMap(foldersMap map[string]*folderInfo) (map[string]interface{}, *goahocorasick.Machine) {
    var excludeFolders []string
    var includedFiles = make(map[string]interface{})
    for k, v := range foldersMap {
        if v.project == nil {
            continue
        }

        // Add project base + exlude subfolder into exclude folders list
        for _, s := range subfolderToExclude {
            sub := filepath.Join(k, s)
            excludeFolders = append(excludeFolders, sub)
        }

        // Add compiles, contents and nones into included files map
        for _, c := range v.project.Compiles {
            fp := filepath.Join(k, c.Path)
            includedFiles[fp] = nil
        }

        if v.project.Contents != nil {
            for _, c := range v.project.Contents {
                fp := filepath.Join(k, c.Path)
                includedFiles[fp] = nil
            }
        }

        if v.project.Nones != nil {
            for _, c := range v.project.Nones {
                fp := filepath.Join(k, c.Path)
                includedFiles[fp] = nil
            }
        }
    }
    excludeFoldersMachine := createAhoCorasickMachine(excludeFolders)
    return includedFiles, excludeFoldersMachine
}
