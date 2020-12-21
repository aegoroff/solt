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
	includedFiles  []string
	nonExistence   bool
	filer          sys.Filer
	lost           Matcher
}

func newLostFilesLogic(nonExistence bool, foundFiles []string, foldersToIgnore c9s.StringHashSet, filer sys.Filer) *lostFilesLogic {
	return &lostFilesLogic{
		foundFiles:     foundFiles,
		excludeFolders: foldersToIgnore,
		unexistFiles:   make(map[string][]string),
		nonExistence:   nonExistence,
		filer:          filer,
	}
}

// initialize fills subfoldersToExclude, excludeFolders, includedFiles and unexistFiles
func (lf *lostFilesLogic) initialize(projects []*msvc.MsbuildProject) error {
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
			lf.excludeFolders.Add(sub)
		}

		// In case of SDK projects all files inside project folder are considered included
		if prj.Project.IsSdkProject() {
			lf.excludeFolders.Add(pdir)
		}

		lf.includedFiles = append(lf.includedFiles, prj.Files()...)

		lf.addToUnexistIfNeeded(prj.Path)
	}

	return lf.initializeLostMatcher()
}

func (lf *lostFilesLogic) addToUnexistIfNeeded(project string) {
	if !lf.nonExistence {
		return
	}

	nonexist := lf.filer.CheckExistence(lf.includedFiles)

	if len(nonexist) > 0 {
		lf.unexistFiles[project] = append(lf.unexistFiles[project], nonexist...)
	}
}

func (lf *lostFilesLogic) initializeLostMatcher() error {
	excludes, err := NewPartialMatcher(lf.excludeFolders.Items(), normalize)
	if err != nil {
		return err
	}

	includes := NewExactMatch(lf.includedFiles)

	lf.lost = NewLostItemMatcher(includes, excludes)

	return nil
}

func (lf *lostFilesLogic) find() []string {
	if lf.lost == nil {
		return []string{}
	}

	var result []string
	for _, file := range lf.foundFiles {
		if lf.lost.Match(file) {
			result = append(result, file)
		}
	}

	return result
}

func (lf *lostFilesLogic) remove(lostFiles []string) {
	lf.filer.Remove(lostFiles)
}
