package lostfiles

import (
	"github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/internal/sys"
	"solt/msvc"
	"strings"
)

type skpper struct {
	folders collections.StringHashSet
}

func newSkipper() *skpper {
	return &skpper{
		folders: make(collections.StringHashSet),
	}
}

func (h *skpper) skipped() []string {
	return h.folders.Items()
}

// Handler executed on each found file in a folder
func (h *skpper) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)

	if strings.EqualFold(ext, msvc.SolutionFileExt) {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		h.folders.Add(ppath)
	}
}

func (h *skpper) fromProject(prj *msvc.MsbuildProject) {
	pdir := filepath.Dir(prj.Path)

	subfolders := []string{"obj"}

	// Exclude output paths too something like bin\Debug, bin\Release etc.
	if prj.Project.OutputPaths != nil {
		subfolders = append(subfolders, prj.Project.OutputPaths...)
	}

	// Add project base + exclude subfolder into exclude folders list
	for _, s := range subfolders {
		sub := filepath.Join(pdir, sys.ToValidPath(s))
		h.folders.Add(sub)
	}

	// In case of SDK projects all files inside project folder are considered included
	if prj.Project.IsSdkProject() {
		h.folders.Add(pdir)
	}
}
