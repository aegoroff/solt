package lostfiles

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/msvc"
)

type enumerator struct {
	excludeFolders c9s.StringHashSet
	includes       []string
	exister        exister
}

func newEnumerator(foldersToIgnore c9s.StringHashSet, exister exister) *enumerator {
	return &enumerator{
		excludeFolders: foldersToIgnore,
		exister:        exister,
	}
}

func (e *enumerator) includedFiles() []string {
	return e.includes
}

func (e *enumerator) excludedFolders() []string {
	return e.excludeFolders.Items()
}

// enumerate fills excludeFolders, includes and validates files existence
func (e *enumerator) enumerate(projects []*msvc.MsbuildProject) {
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
			e.excludeFolders.Add(sub)
		}

		// In case of SDK projects all files inside project folder are considered included
		if prj.Project.IsSdkProject() {
			e.excludeFolders.Add(pdir)
		}

		includes := prj.Files()
		e.includes = append(e.includes, includes...)

		e.exister.exist(prj.Path, includes)
	}
}
