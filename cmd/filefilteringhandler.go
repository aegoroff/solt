package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/msvc"
	"strings"
)

type fileFilteringHandler struct {
	foundFiles      []string
	foldersToIgnore c9s.StringHashSet
	filter          string
}

func newFileFilteringHandler(filter string) *fileFilteringHandler {
	return &fileFilteringHandler{
		foundFiles:      make([]string, 0),
		foldersToIgnore: make(c9s.StringHashSet),
		filter:          filter,
	}
}

// Handler executed on each found file in a folder
func (h *fileFilteringHandler) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)
	if strings.EqualFold(ext, h.filter) {
		h.foundFiles = append(h.foundFiles, path)
	}

	if strings.EqualFold(ext, msvc.SolutionFileExt) {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		h.foldersToIgnore.Add(ppath)
	}
}
