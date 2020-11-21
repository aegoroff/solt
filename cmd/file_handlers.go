package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/msvc"
	"strings"
)

type fileExtFilteringHandler struct {
	foundFiles []string
	filter     string
}

type folderIgnoreHandler struct {
	folders c9s.StringHashSet
}

func newFileExtFilteringHandler(filter string) *fileExtFilteringHandler {
	return &fileExtFilteringHandler{
		foundFiles: make([]string, 0),
		filter:     filter,
	}
}

func newFolderIgnoreHandler() *folderIgnoreHandler {
	return &folderIgnoreHandler{
		folders: make(c9s.StringHashSet),
	}
}

// Handler executed on each found file in a folder
func (h *fileExtFilteringHandler) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)
	if strings.EqualFold(ext, h.filter) {
		h.foundFiles = append(h.foundFiles, path)
	}
}

// Handler executed on each found file in a folder
func (h *folderIgnoreHandler) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)

	if strings.EqualFold(ext, msvc.SolutionFileExt) {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		h.folders.Add(ppath)
	}
}
