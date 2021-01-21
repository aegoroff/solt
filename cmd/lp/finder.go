package lp

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"os"
	"path/filepath"
	"solt/internal/fw"
	"solt/msvc"
	"strings"
)

type finder struct {
	allFilesPaths c9s.StringHashSet
}

func newFinder() *finder {
	return &finder{
		allFilesPaths: c9s.NewStringHashSet(),
	}
}

func (f *finder) filter(all []*msvc.MsbuildProject, withinSolution []string) ([]string, []string) {
	// Create projects matching machine
	within := fw.NewExactMatch(withinSolution)

	n := 0
	for _, p := range all {
		if within.Match(p.Path()) {
			f.selectFilePaths(p)
		} else {
			all[n] = p
			n++
		}
	}

	all = all[:n]
	return f.separate(all)
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

	ldMatch, err := fw.NewPartialMatcher(lostDirs, strings.ToUpper)
	if err == nil {
		for fp := range f.allFilesPaths {
			filePaths.Add(fp)
			r := ldMatch.Search(fp)
			filePaths.AddRange(r...)
		}
	}

	m, err := fw.NewPartialMatcher(filePaths.Items(), strings.ToUpper)
	if err != nil {
		return fw.NewMatchNothing()
	}
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
