package lostfiles

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/cmd/api"
	"solt/msvc"
	"strings"
)

type lostFilesLogic struct {
	foundFiles     []string
	excludeFolders c9s.StringHashSet
	includedFiles  []string
	nonExistence   bool
	exister        *api.Exister
	lost           api.Matcher
}

func newLostFilesLogic(nonExistence bool, foundFiles []string, foldersToIgnore c9s.StringHashSet, exister *api.Exister) *lostFilesLogic {
	return &lostFilesLogic{
		foundFiles:     foundFiles,
		excludeFolders: foldersToIgnore,
		nonExistence:   nonExistence,
		exister:        exister,
	}
}

// initialize fills subfoldersToExclude, excludeFolders, includedFiles and unexistFiles
func (lf *lostFilesLogic) initialize(projects []*msvc.MsbuildProject) error {
	subfoldersToExclude := []string{"obj"}

	for _, prj := range projects {
		pdir := filepath.Dir(prj.Path)

		// Exclude output paths too something like bin\Debug, bin\Release etc.
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

		includes := prj.Files()
		lf.includedFiles = append(lf.includedFiles, includes...)

		lf.validateExistence(prj.Path, includes)
	}

	return lf.initializeLostMatcher()
}

func (lf *lostFilesLogic) validateExistence(project string, includes []string) {
	if lf.nonExistence {
		lf.exister.Validate(project, includes)
	}
}

func (lf *lostFilesLogic) initializeLostMatcher() error {
	excludes, err := api.NewPartialMatcher(lf.excludeFolders.Items(), strings.ToUpper)
	if err != nil {
		return err
	}

	includes := api.NewExactMatch(lf.includedFiles)

	lf.lost = api.NewLostItemMatcher(includes, excludes)

	return nil
}

func (lf *lostFilesLogic) find() []string {
	return api.Filter(lf.foundFiles, lf.lost)
}
