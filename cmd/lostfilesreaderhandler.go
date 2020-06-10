package cmd

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"path/filepath"
	"solt/internal/sys"
	"solt/msvc"
	"strings"
)

type lostFilesHandler struct {
	fs                  afero.Fs
	foundFiles          []string
	excludeFolders      collections.StringHashSet
	lostFilesFilter     string
	unexistFiles        map[string][]string
	includedFiles       collections.StringHashSet
	subfoldersToExclude []string
	nonExistence        bool
	filer               sys.Filer
}

func newLostFilesHandler(lostFilesFilter string, nonExistence bool, fs afero.Fs) *lostFilesHandler {
	return &lostFilesHandler{
		fs:                  fs,
		foundFiles:          make([]string, 0),
		excludeFolders:      make(collections.StringHashSet),
		lostFilesFilter:     lostFilesFilter,
		unexistFiles:        make(map[string][]string),
		includedFiles:       make(collections.StringHashSet),
		subfoldersToExclude: []string{"obj"},
		nonExistence:        nonExistence,
		filer:               sys.NewFiler(fs, appWriter),
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

		// Exclude output paths too
		if prj.Project.OutputPaths != nil {
			r.subfoldersToExclude = append(r.subfoldersToExclude, prj.Project.OutputPaths...)
		}

		// Add project base + exclude subfolder into exclude folders list
		for _, s := range r.subfoldersToExclude {
			sub := filepath.Join(pdir, s)
			r.excludeFolders.Add(sub)
		}

		// In case of SDK projects all files inside project folder are considered included
		if prj.Project.IsSdkProject() {
			r.excludeFolders.Add(pdir)
		}

		// Add compiles, contents and nones into included files map
		includes := msvc.GetFilesIncludedIntoProject(prj)
		for _, f := range includes {
			r.includedFiles.Add(normalize(f))
		}

		r.checkExistence(prj.Path, includes)
	}
}

func (r *lostFilesHandler) checkExistence(project string, includes []string) {
	if !r.nonExistence {
		return
	}

	nonexist := r.filer.CheckExistence(includes)

	if len(nonexist) > 0 {
		r.unexistFiles[project] = append(r.unexistFiles[project], nonexist...)
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

func (r *lostFilesHandler) removeLostFiles(lostFiles []string) {
	r.filer.Remove(lostFiles)
}
