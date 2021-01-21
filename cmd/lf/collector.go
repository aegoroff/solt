package lf

import (
	"path/filepath"
	"strings"
)

type collector struct {
	files  []string
	filter string
}

// newCollector creates new collector instance
// filter - file extension to collect files that match it
func newCollector(filter string) *collector {
	return &collector{
		files:  make([]string, 0),
		filter: filter,
	}
}

// Handler executed on each found file in a folder
func (h *collector) Handler(path string) {
	// Add file to filtered files slice
	ext := filepath.Ext(path)
	if strings.EqualFold(ext, h.filter) {
		h.files = append(h.files, path)
	}
}
