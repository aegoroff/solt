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
	lostMatcher    Matcher
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
func (h *lostFilesLogic) initialize(projects []*msvc.MsbuildProject) error {
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

	return h.initializeLostMatcher()
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

func (h lostFilesLogic) initializeLostMatcher() error {
	excludes, err := NewPartialMatcher(h.excludeFolders.ItemsDecorated(normalize))
	if err != nil {
		return err
	}

	includes := NewExactMatchHS(&h.includedFiles)

	h.lostMatcher = NewLostItemMatcher(includes, excludes, normalize)

	return nil
}

func (h *lostFilesLogic) find() []string {
	if h.lostMatcher == nil {
		return []string{}
	}

	var result []string
	for _, file := range h.foundFiles {
		if h.lostMatcher.Match(file) {
			result = append(result, file)
		}
	}

	return result
}

func (h *lostFilesLogic) remove(lostFiles []string) {
	h.filer.Remove(lostFiles)
}
