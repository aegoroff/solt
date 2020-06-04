package cmd

import (
	"github.com/aegoroff/godatastruct/collections"
	"path/filepath"
	"solt/internal/msvc"
)

type lostFilesHandler struct {
	foundFiles      []string
	excludeFolders  collections.StringHashSet
	lostFilesFilter string
}

func newLostFilesHandler(lostFilesFilter string) *lostFilesHandler {
	return &lostFilesHandler{
		foundFiles:      make([]string, 0),
		excludeFolders:  make(collections.StringHashSet),
		lostFilesFilter: lostFilesFilter,
	}
}

func (r *lostFilesHandler) Handler(path string) {
	ef := normalize(r.lostFilesFilter)
	sln := normalize(msvc.SolutionFileExt)

	// Add file to filtered files slice
	ext := normalize(filepath.Ext(path))
	if ext == ef {
		r.foundFiles = append(r.foundFiles, path)
	}

	if ext == sln {
		dir, _ := filepath.Split(path)
		ppath := filepath.Join(dir, "packages")
		r.excludeFolders.Add(ppath)
	}
}
