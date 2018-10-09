package main

import (
	"fmt"
	"os"
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

	lostFiles, unexistFiles := findLostFiles(folders, packagesFolders, foundFiles)

	sortAndOutput(lostFiles)

	if len(unexistFiles) > 0 {
		fmt.Printf("\nThese files included into projects but not exist in the file system.\n")
	}

	sortAndOutput(unexistFiles)

	return nil
}

func findLostFiles(folders []*folderInfo, packagesFolders map[string]interface{}, foundFiles []string) ([]string, []string) {
	includedFiles, excludedFolders, unexistFiles := createIncludedFilesAndExcludedFolders(folders)
	for k := range packagesFolders {
		excludedFolders = append(excludedFolders, k)
	}

	exmach := createAhoCorasickMachine(excludedFolders)
	var result []string
	for _, file := range foundFiles {
		if _, ok := includedFiles[strings.ToUpper(file)]; !ok && !Match(exmach, file) {
			result = append(result, file)
		}
	}

	return result, unexistFiles
}

func createIncludedFilesAndExcludedFolders(folders []*folderInfo) (map[string]interface{}, []string, []string) {
	var excludeFolders []string
	var unexistFiles []string
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
			if _, err := os.Stat(f); os.IsNotExist(err) {
				unexistFiles = append(unexistFiles, f)
			}
		}
	}
	return includedFiles, excludeFolders, unexistFiles
}
