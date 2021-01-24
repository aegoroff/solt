package va

import (
	"github.com/spf13/afero"
	"solt/internal/fw"
	"solt/msvc"
	"sort"
)

type validator struct {
	fs          afero.Fs
	sourcesPath string
	act         actioner
	sdk         *sdkProjects
	tt          *totals
}

func newValidator(fs afero.Fs, sourcesPath string, act actioner) *validator {
	return &validator{
		fs:          fs,
		sourcesPath: sourcesPath,
		act:         act,
		tt:          &totals{},
	}
}

func (va *validator) validate() {
	foldersTree := msvc.ReadSolutionDir(va.sourcesPath, va.fs)

	sols, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)
	va.sdk = newSdkProjects(allProjects)
	va.tt.solutions = int64(len(sols))

	solutions := fw.SolutionSlice(sols)
	sort.Sort(solutions)

	solutions.Foreach(va)
}

func (va *validator) Solution(sol *msvc.VisualStudioSolution) {
	search := newSdkSearcher(va.sdk, sol)
	it := msvc.NewProjectIterator(sol, search)
	gr := newGraph(it)
	find := newFinder(gr)
	redundants := find.findAll()

	va.tt.projects += gr.nextID - 1
	if len(redundants) > 0 {
		va.tt.problemSolutions++
		va.tt.problemProjects += int64(len(redundants))
		for _, set := range redundants {
			va.tt.redundantRefs += int64(set.Count())
		}
	}

	va.act.action(sol.Path(), redundants)
}
