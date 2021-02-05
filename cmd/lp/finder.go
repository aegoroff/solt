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
	within := fw.NewExactMatch(withinSolution, true)
	lost := make([]string, len(all))

	n := 0
	for _, p := range all {
		pp := p.Path()
		if within.Match(pp) {
			f.selectFilePaths(p)
		} else {
			lost[n] = pp
			n++
		}
	}

	return f.separate(lost[:n])
}

func (f *finder) separate(lost []string) ([]string, []string) {
	lostWithIncludes := make([]string, len(lost))

	filesFoldersM, lostDirs := f.newMatcher(lost)

	i := 0
	j := 0
	for ix, lp := range lost {
		d := lostDirs[ix]
		if filesFoldersM.Match(d) {
			lostWithIncludes[j] = lp
			j++
		} else {
			lost[i] = lp
			i++
		}
	}
	return lost[:i], lostWithIncludes[:j]
}

func (f *finder) newMatcher(allLost []string) (fw.Matcher, []string) {
	filePaths := make(c9s.StringHashSet, f.allFilesPaths.Count())
	lostDirs := make([]string, len(allLost))

	for i, lp := range allLost {
		lostDirs[i] = dir(lp)
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
		return fw.NewMatchNothing(), lostDirs
	}
	return m, lostDirs
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
