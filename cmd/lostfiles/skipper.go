package lostfiles

import (
	"github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/internal/sys"
	"solt/msvc"
	"strings"
)

// skipper is just collector that collects folders
// which files have to be skipped
type skipper struct {
	items collections.StringHashSet
}

func newSkipper() *skipper {
	return &skipper{
		items: make(collections.StringHashSet),
	}
}

func (h *skipper) folders() []string {
	return h.items.Items()
}

// Handler executed on each found file in a folder
func (h *skipper) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)

	if strings.EqualFold(ext, msvc.SolutionFileExt) {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		h.items.Add(ppath)
	}
}

func (h *skipper) fromProject(prj *msvc.MsbuildProject) {
	pdir := filepath.Dir(prj.Path())

	// In case of SDK projects all files inside project folder are considered included
	if prj.Project.IsSdkProject() {
		h.items.Add(pdir)
		return
	}

	subfolders := []string{"obj"}

	// Exclude output paths too something like bin\Debug, bin\Release etc.
	if prj.Project.OutputPaths != nil {
		subfolders = append(subfolders, prj.Project.OutputPaths...)
	}

	// Add project base + exclude subfolder into exclude folders list
	for _, s := range subfolders {
		sub := filepath.Join(pdir, sys.ToValidPath(s))
		h.items.Add(sub)
	}
}
