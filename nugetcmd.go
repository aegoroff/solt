package main

import (
    "fmt"
)

func nugetcmd(opt options) error {

    foldersMap := readProjectDir(opt.Path, func(we *walkEntry) {})

    for k, v := range foldersMap {
        if v.packages == nil {
            continue
        }
        fmt.Printf("%s\n", k)
        for _, p := range v.packages.Packages {
            fmt.Printf("\t%s %s\n", p.Id, p.Version)
        }
    }

    return nil
}
