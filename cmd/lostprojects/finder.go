package lostprojects

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"os"
	"path/filepath"
	"solt/cmd/fw"
	"solt/msvc"
	"strings"
)

type finder struct {
	allFilesPaths c9s.StringHashSet
}

func newFinder() *finder {
	return &finder{
		allFilesPaths: make(c9s.StringHashSet),
	}
}

func (f *finder) filter(all []*msvc.MsbuildProject, withinSolution []string) ([]string, []string) {
	// Create projects matching machine
	within := fw.NewExactMatch(withinSolution)

	lost := all[:0]

	for _, p := range all {
		if within.Match(p.Path()) {
			f.selectFilePaths(p)
		} else {
			lost = append(lost, p)
		}
	}

	return f.separate(lost)
}

func (f *finder) separate(allLost []*msvc.MsbuildProject) ([]string, []string) {
	lost := make([]string, 0, len(allLost))
	var lostWithIncludes []string

	lostDirs := make([]string, 0, len(allLost))
	for _, lp := range allLost {
		lostDirs = append(lostDirs, dir(lp.Path()))
	}

	lostDirMatch, _ := fw.NewPartialMatcher(lostDirs, strings.ToUpper)

	allFilesPaths := f.filesPaths()
	for _, lp := range allLost {
		if lostDirMatch.Match(lp.Path()) && fw.MatchAny(allFilesPaths, lostDirMatch) {
			lostWithIncludes = append(lostWithIncludes, lp.Path())
		} else {
			lost = append(lost, lp.Path())
		}
	}
	return lost, lostWithIncludes
}

func (f *finder) filesPaths() []string {
	return f.allFilesPaths.ItemsDecorated(trailPathSeparator)
}

func (f *finder) selectFilePaths(p *msvc.MsbuildProject) {
	for _, s := range p.Items() {
		f.allFilesPaths.Add(dir(s))
	}
}

func dir(path string) string {
	return trailPathSeparator(filepath.Dir(path))
}

func trailPathSeparator(s string) string {
	return s + string(os.PathSeparator)
}
