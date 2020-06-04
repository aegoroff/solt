package cmd

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"solt/internal/msvc"
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

func (r *lostFilesHandler) Handler(path string) {
	ef := normalize(r.lostFilesFilter)
	sln := normalize(msvc.SolutionFileExt)

	// Add file to filtered files slice
	ext := normalize(filepath.Ext(path))
	if ext == ef {
		r.foundFiles = append(r.foundFiles, path)
	}

	if ext == sln {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		r.excludeFolders.Add(ppath)
	}
}

func (r *lostFilesHandler) projectHandler(prj *msvc.MsbuildProject, fo *msvc.Folder) {
	// Add project base + exclude subfolder into exclude folders list
	for _, s := range subfolderToExclude {
		sub := filepath.Join(fo.Path, s)
		r.excludeFolders.Add(sub)
	}

	// Exclude output paths too
	if prj.Project.OutputPaths != nil {
		for _, out := range prj.Project.OutputPaths {
			sub := filepath.Join(fo.Path, out)
			r.excludeFolders.Add(sub)
		}
	}

	// In case of SDK projects all files inside project folder are considered included
	if prj.Project.IsSdkProject() {
		r.excludeFolders.Add(filepath.Dir(prj.Path))
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
