package lostfiles

import (
	"path/filepath"
	"strings"
)

type fileCollector struct {
	files  []string
	filter string
}

// newFileCollector creates new collector instance
// filter - file extension to collect files that match it
func newFileCollector(filter string) *fileCollector {
	return &fileCollector{
		files:  make([]string, 0),
		filter: filter,
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
