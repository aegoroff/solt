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
	var lost []string
	var lostWithIncludes []string

	pathMatch, _ := fw.NewPartialMatcher(f.filesPaths(), strings.ToUpper)

	for _, lp := range allLost {
		if pathMatch.Match(lp.Path()) {
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
		d := filepath.Dir(s)
		f.allFilesPaths.Add(d)
	}
}

func trailPathSeparator(s string) string {
	return s + string(os.PathSeparator)
}
