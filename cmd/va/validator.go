package va

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/spf13/afero"
	"solt/internal/fw"
	"solt/msvc"
	"solt/msvc/graph"
	"sort"
)

type validator struct {
	fs              afero.Fs
	sourcesPath     string
	act             actioner
	sdk             *sdkProjects
	tt              *totals
	problemProjects c9s.StringHashSet
}

func newValidator(fs afero.Fs, sourcesPath string, act actioner) *validator {
	return &validator{
		fs:              fs,
		sourcesPath:     sourcesPath,
		act:             act,
		tt:              &totals{},
		problemProjects: c9s.NewStringHashSet(),
	}
}

func (va *validator) validate() {
	foldersTree := msvc.ReadSolutionDir(va.sourcesPath, va.fs)

	sols, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)
	va.sdk = newSdkProjects(allProjects)
	va.tt.solutions = int64(len(sols))
	va.tt.projects = va.sdk.count()

	solutions := fw.SolutionSlice(sols)
	sort.Sort(solutions)

	solutions.Foreach(va)
}

func (va *validator) Solution(sol *msvc.VisualStudioSolution) {
	search := newSdkSearcher(va.sdk, sol)
	it := msvc.NewProjectIterator(sol, search)
	gr := graph.New(it)
	find := newFinder(gr)
	redundants := find.findAll()

	if len(redundants) > 0 {
		va.tt.problemSolutions++
		for prj, set := range redundants {
			if va.problemProjects.Contains(prj) {
				continue
			}
			va.problemProjects.Add(prj)
			va.tt.problemProjects++
			va.tt.redundantRefs += int64(set.Count())
		}
	}

	va.act.action(sol.Path(), redundants)
}
