package va

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"solt/internal/fw"
	"solt/msvc"
	"sort"
)

type validator struct {
	fs          afero.Fs
	sourcesPath string
	act         actioner
	sdkProjects rbtree.RbTree
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
	va.tt.solutions = int64(len(sols))

	va.onlySdkProjects(allProjects)
	va.tt.projects = va.sdkProjects.Len()

	solutions := fw.SolutionSlice(sols)
	sort.Sort(solutions)

	solutions.Foreach(va)
}

func (va *validator) onlySdkProjects(allProjects []*msvc.MsbuildProject) {
	va.sdkProjects = rbtree.New()

	for _, p := range allProjects {
		if p.IsSdkProject() {
			va.sdkProjects.Insert(p)
		}
	}
}

func (va *validator) Solution(sol *msvc.VisualStudioSolution) {
	gr := newGraph(sol, va.sdkProjects)
	redundants := findRedundants(gr)

	if len(redundants) > 0 {
		va.tt.problemSolutions++
		va.tt.problemProjects += int64(len(redundants))
		for _, set := range redundants {
			va.tt.redundantRefs += int64(set.Count())
		}
	}

	va.act.action(sol.Path(), redundants)
}

func findRedundants(g *graph) map[string]c9s.StringHashSet {
	result := make(map[string]c9s.StringHashSet)
	find := newFinder(g.allPaths())

	g.foreach(func(n *node) {
		found, ok := find.find(n.refs)

		if ok {
			result[n.String()] = found
		}
	})

	return result
}
