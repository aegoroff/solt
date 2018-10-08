package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

var subfolderToExclude = []string{
	"obj",
}

func lostfilescmd(opt options) error {

	filter := csharpCodeFileExt
	if len(opt.LostFiles.Filter) > 0 {
		filter = opt.LostFiles.Filter
	}

	var foundFiles []string
	var packagesFolders = make(map[string]interface{})
	folders := readProjectDir(opt.Path, func(we *walkEntry) {
		// Add file to filtered files slice
		ext := strings.ToLower(filepath.Ext(we.Name))
		if ext == filter {
			fp := filepath.Join(we.Parent, we.Name)
			foundFiles = append(foundFiles, fp)
		}

		if ext == solutionFileExt {
			ppath := filepath.Join(we.Parent, "packages")
			if _, ok := packagesFolders[ppath]; !ok {
				packagesFolders[ppath] = nil
			}
		}
	})

	includedFiles, excludedFolders := createIncludedFilesAndExcludedFolders(folders)

	for k := range packagesFolders {
		excludedFolders = append(excludedFolders, k)
	}

	findLostFiles(excludedFolders, foundFiles, includedFiles)

	return nil
}

func findLostFiles(excludedFolders []string, foundFiles []string, includedFiles map[string]interface{}) {
	exmach := createAhoCorasickMachine(excludedFolders)
	for _, file := range foundFiles {
		if _, ok := includedFiles[strings.ToUpper(file)]; !ok && !Match(exmach, file) {
			fmt.Println(file)
		}
	}
}

func createIncludedFilesAndExcludedFolders(folders []*folderInfo) (map[string]interface{}, []string) {
	var excludeFolders []string
	var includedFiles = make(map[string]interface{})
	for _, info := range folders {
		if info.project == nil {
			continue
		}

		// Add project base + exclude subfolder into exclude folders list
		parent := filepath.Dir(*info.projectPath)
		for _, s := range subfolderToExclude {
			sub := filepath.Join(parent, s)
			excludeFolders = append(excludeFolders, sub)
		}

		if info.project.OutputPaths != nil {
			for _, s := range info.project.OutputPaths {
				sub := filepath.Join(parent, s)
				excludeFolders = append(excludeFolders, sub)
			}
		}

		// Add compiles, contents and nones into included files map
		filesIncluded := getFilesIncludedIntoProject(info)
		for _, f := range filesIncluded {
			includedFiles[strings.ToUpper(f)] = nil
		}
	}
	return includedFiles, excludeFolders
}
