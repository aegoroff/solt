package lf

import (
	"path/filepath"
	"solt/internal/sys"
	"solt/msvc"
	"strings"

	c9s "github.com/aegoroff/godatastruct/collections"
)

// skipper is just collector that collects folders
// which files have to be skipped
type skipper struct {
	items c9s.HashSet[string]
}

func newSkipper() *skipper {
	return &skipper{
		items: c9s.NewHashSet[string](),
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
	if prj.IsSdkProject() {
		h.items.Add(pdir)
		return
	}

	var subfolders []string
	const objPath = "obj"
	// Exclude output paths too something like bin\Debug, bin\Release etc.
	if prj.Project.OutputPaths != nil {
		subfolders = make([]string, len(prj.Project.OutputPaths)+1)
		i := copy(subfolders, prj.Project.OutputPaths)
		subfolders[i] = objPath
	} else {
		subfolders = []string{objPath}
	}

	// Add project base + exclude subfolder into exclude folders list
	for _, s := range subfolders {
		sub := filepath.Join(pdir, sys.ToValidPath(s))
		h.items.Add(sub)
	}
}
