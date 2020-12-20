package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/internal/sys"
	"solt/msvc"
)

type lostFilesLogic struct {
	foundFiles     []string
	excludeFolders c9s.StringHashSet
	unexistFiles   map[string][]string
	includedFiles  c9s.StringHashSet
	nonExistence   bool
	filer          sys.Filer
}

func newLostFilesLogic(nonExistence bool, foundFiles []string, foldersToIgnore c9s.StringHashSet, filer sys.Filer) *lostFilesLogic {
	return &lostFilesLogic{
		foundFiles:     foundFiles,
		excludeFolders: foldersToIgnore,
		unexistFiles:   make(map[string][]string),
		includedFiles:  make(c9s.StringHashSet),
		nonExistence:   nonExistence,
		filer:          filer,
	}
}

// initialize fills subfoldersToExclude, excludeFolders, includedFiles and unexistFiles
func (h *lostFilesLogic) initialize(projects []*msvc.MsbuildProject) {
	subfoldersToExclude := []string{"obj"}

	for _, prj := range projects {
		pdir := filepath.Dir(prj.Path)

		// Exclude output paths too
		if prj.Project.OutputPaths != nil {
			subfoldersToExclude = append(subfoldersToExclude, prj.Project.OutputPaths...)
		}

		// Add project base + exclude subfolder into exclude folders list
		for _, s := range subfoldersToExclude {
			sub := filepath.Join(pdir, s)
			h.excludeFolders.Add(sub)
		}

		// In case of SDK projects all files inside project folder are considered included
		if prj.Project.IsSdkProject() {
			h.excludeFolders.Add(pdir)
		}

		// Add compiles, contents and nones into included files map
		includes := prj.Files()
		for _, f := range includes {
			h.includedFiles.Add(normalize(f))
		}

		h.addToUnexistIfNeeded(prj.Path, includes)
	}
}

func (h *lostFilesLogic) addToUnexistIfNeeded(project string, includes []string) {
	if !h.nonExistence {
		return
	}

	nonexist := h.filer.CheckExistence(includes)

	if len(nonexist) > 0 {
		h.unexistFiles[project] = append(h.unexistFiles[project], nonexist...)
	}
}

func (h *lostFilesLogic) find() ([]string, error) {
	excludes, err := NewPartialMatcher(h.excludeFolders.ItemsDecorated(normalize))
	if err != nil {
		return nil, err
	}

	includes := NewExactMatchHS(&h.includedFiles)

	lostMatch := NewLostItemMatcher(includes, excludes, normalize)

	var result []string
	for _, file := range h.foundFiles {
		if lostMatch.Match(file) {
			result = append(result, file)
		}
	}

	return result, err
}

func (h *lostFilesLogic) remove(lostFiles []string) {
	h.filer.Remove(lostFiles)
}
