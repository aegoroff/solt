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

	filesFoldersM := f.newMatcher(allLost)

	for _, lp := range allLost {
		d := dir(lp.Path())
		if filesFoldersM.Match(d) {
			lostWithIncludes = append(lostWithIncludes, lp.Path())
		} else {
			lost = append(lost, lp.Path())
		}
	}
	return lost, lostWithIncludes
}

func (f *finder) newMatcher(allLost []*msvc.MsbuildProject) fw.Matcher {
	filePaths := make(c9s.StringHashSet, len(f.allFilesPaths))
	lostDirs := make([]string, len(allLost))

	for i, lp := range allLost {
		lostDirs[i] = dir(lp.Path())
	}

	dm, err := fw.NewPartialMatcher(lostDirs, strings.ToUpper)
	if err == nil {
		for path := range f.allFilesPaths {
			filePaths.Add(path)
			r := dm.Search(path)
			filePaths.AddRange(r...)
		}
	}

	m, _ := fw.NewPartialMatcher(filePaths.Items(), strings.ToUpper)
	return m
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