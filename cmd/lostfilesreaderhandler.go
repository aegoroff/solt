package cmd

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"solt/internal/msvc"
	"strings"
)

type lostFilesHandler struct {
	fs              afero.Fs
	foundFiles      []string
	excludeFolders  collections.StringHashSet
	lostFilesFilter string
	unexistFiles    map[string][]string
	includedFiles   collections.StringHashSet
}

func newLostFilesHandler(lostFilesFilter string, fs afero.Fs) *lostFilesHandler {
	return &lostFilesHandler{
		fs:              fs,
		foundFiles:      make([]string, 0),
		excludeFolders:  make(collections.StringHashSet),
		lostFilesFilter: lostFilesFilter,
		unexistFiles:    make(map[string][]string),
		includedFiles:   make(collections.StringHashSet),
	}
}

// executed on each found file in a folder
func (r *lostFilesHandler) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)
	if strings.EqualFold(ext, r.lostFilesFilter) {
		r.foundFiles = append(r.foundFiles, path)
	}

	if strings.EqualFold(ext, msvc.SolutionFileExt) {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		r.excludeFolders.Add(ppath)
	}
}

// Executed on each found folder that contains msbuild projects
func (r *lostFilesHandler) projectHandler(projects []*msvc.MsbuildProject) {
	for _, prj := range projects {
		pdir := filepath.Dir(prj.Path)
		// Add project base + exclude subfolder into exclude folders list
		for _, s := range subfolderToExclude {
			sub := filepath.Join(pdir, s)
			r.excludeFolders.Add(sub)
		}

		// Exclude output paths too
		if prj.Project.OutputPaths != nil {
			for _, out := range prj.Project.OutputPaths {
				sub := filepath.Join(pdir, out)
				r.excludeFolders.Add(sub)
			}
		}

		// In case of SDK projects all files inside project folder are considered included
		if prj.Project.IsSdkProject() {
			r.excludeFolders.Add(pdir)
		}

		// Add compiles, contents and nones into included files map
		filesIncluded := msvc.GetFilesIncludedIntoProject(prj)
		for _, f := range filesIncluded {
			normalized := normalize(f)
			r.includedFiles.Add(normalized)
			if _, err := r.fs.Stat(f); os.IsNotExist(err) {
				if found, ok := r.unexistFiles[prj.Path]; ok {
					found = append(found, f)
					r.unexistFiles[prj.Path] = found
				} else {
					r.unexistFiles[prj.Path] = []string{f}
				}
			}
		}
	}
}

func (r *lostFilesHandler) findLostFiles() ([]string, error) {
	excludes, err := NewPartialMatcher(r.excludeFolders.ItemsDecorated(normalize))
	if err != nil {
		return nil, err
	}

	includes := NewExactMatchHS(&r.includedFiles)

	var result []string
	for _, file := range r.foundFiles {
		normalized := normalize(file)
		if !includes.Match(normalized) && !excludes.Match(normalized) {
			result = append(result, file)
		}
	}

	return result, err
}
