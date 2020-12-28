package lostfiles

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/msvc"
	"strings"
)

type fileCollector struct {
	files  []string
	filter string
}

type ignoredFoldersCollector struct {
	folders c9s.StringHashSet
}

// newFileCollector creates new collector instance
// filter - file extension to collect files that match it
func newFileCollector(filter string) *fileCollector {
	return &fileCollector{
		files:  make([]string, 0),
		filter: filter,
	}
}

func newFoldersCollector() *ignoredFoldersCollector {
	return &ignoredFoldersCollector{
		folders: make(c9s.StringHashSet),
	}
}

// Handler executed on each found file in a folder
func (h *fileCollector) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)
	if strings.EqualFold(ext, h.filter) {
		h.files = append(h.files, path)
	}
}

// Handler executed on each found file in a folder
func (h *ignoredFoldersCollector) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)

	if strings.EqualFold(ext, msvc.SolutionFileExt) {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		h.folders.Add(ppath)
	}
}
