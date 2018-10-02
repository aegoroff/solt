package main

import (
    "fmt"
    "path/filepath"
    "strings"
)

var excludeSubFolders = []string{
    "obj",
    "bin",
}

func lostcmd(opt options) error {
    var includedFiles = make(map[string]interface{})
    var filteredFiles []string

    filter := CSharpCodeFileExt
    if len(opt.Lost.Filter) > 0 {
        filter = opt.Lost.Filter
    }

    var excludeFolders []string

    foldersMap := readProjectDir(opt.Path, func(we *walkEntry) {
        // Add file to filtered files slice
        ext := strings.ToLower(filepath.Ext(we.Name))
        if ext == filter {
            fp := filepath.Join(we.Parent, we.Name)
            filteredFiles = append(filteredFiles, fp)
        }
    })

    for k, v := range foldersMap {
        if v.project == nil {
            continue
        }

        // Add project base + exlude subfolder into exclude folders list
        for _, s := range excludeSubFolders {
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

    for _, file := range filteredFiles {
        if _, ok := includedFiles[file]; !ok {

            if Match(excludeFoldersMachine, file) {
                continue
            }

            fmt.Println(file)
        }
    }

    return nil
}
