package lostprojects

import (
	"solt/cmd/fw"
	"solt/msvc"
)

type finder struct {
	allFiles *fw.Includer
}

func newFinder() *finder {
	return &finder{
		allFiles: fw.NewIncluder(fw.NewNullExister()),
	}
}

func (f *finder) filter(all []*msvc.MsbuildProject, withinSolution []string) ([]string, []string) {
	// Create projects matching machine
	within := fw.NewExactMatch(withinSolution)

	lost := all[:0]

	for _, p := range all {
		if within.Match(p.Path()) {
			f.allFiles.From(p)
		} else {
			lost = append(lost, p)
		}
	}

	return f.separate(lost)
}

func (f *finder) separate(allLost []*msvc.MsbuildProject) ([]string, []string) {
	var lost []string
	var lostWithIncludes []string
	solutionFiles := fw.NewExactMatch(f.allFiles.Includes())

	for _, lp := range allLost {
		if fw.MatchAny(lp.Items(), solutionFiles) {
			lostWithIncludes = append(lostWithIncludes, lp.Path())
		} else {
			lost = append(lost, lp.Path())
		}
	}
	return lost, lostWithIncludes
}
