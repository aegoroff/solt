package va

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/internal/fw"
	"solt/msvc"
	"solt/solution"
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

	for _, project := range allProjects {
		if !project.Project.IsSdkProject() {
			continue
		}
		va.sdkProjects.Insert(project)
	}
}

func (va *validator) Solution(sol *msvc.VisualStudioSolution) {
	g := va.newSolutionGraph(sol)

	redundants := findRedundants(g)
	if len(redundants) > 0 {
		va.tt.problemSolutions++
		va.tt.problemProjects += int64(len(redundants))
		for _, set := range redundants {
			va.tt.redundantRefs += int64(set.Count())
		}
	}

	va.act.action(sol.Path(), redundants)
}

func (va *validator) newSolutionGraph(sln *msvc.VisualStudioSolution) *simple.DirectedGraph {
	solutionPath := filepath.Dir(sln.Path())
	g := simple.NewDirectedGraph()
	nodes := rbtree.New()
	ix := int64(1)
	for _, prj := range sln.Solution.Projects {
		if prj.TypeID == solution.IDSolutionFolder {
			continue
		}

		p := msvc.NewMsbuildProject(filepath.Join(solutionPath, prj.Path))

		msbuild, ok := va.sdkProjects.Search(p)
		if !ok {
			continue
		}

		n := newNode(ix, msbuild.(*msvc.MsbuildProject))
		nodes.Insert(n)
		ix++
		g.AddNode(n)
	}

	createGraphEdges(g, nodes)

	return g
}

func createGraphEdges(g *simple.DirectedGraph, nodes rbtree.RbTree) {
	gn := g.Nodes()

	for gn.Next() {
		to := gn.Node().(*node)
		to.refs = getReferences(to, nodes)
		for _, ref := range to.refs {
			e := g.NewEdge(ref, to)
			g.SetEdge(e)
		}
	}
}

func findRedundants(g *simple.DirectedGraph) map[string]c9s.StringHashSet {
	allPaths := path.DijkstraAllPaths(g)
	result := make(map[string]c9s.StringHashSet)

	gn := g.Nodes()

	for gn.Next() {
		project := gn.Node().(*node)

		rrs := c9s.NewStringHashSet()

		// If there are any paths between project refs
		// they path start point considered redundant link
		allPairs(project.refs, func(from *node, to *node) bool {
			paths, _ := allPaths.AllBetween(from.ID(), to.ID())
			ok := len(paths) > 0
			if ok {
				rrs.Add(from.String())
			}
			// if found iteration will stop
			return ok
		})

		if rrs.Count() > 0 {
			result[project.String()] = rrs
		}
	}

	return result
}

func allPairs(nodes []*node, hasPath func(*node, *node) bool) {
	for _, from := range nodes {
		for _, to := range nodes {
			if from.ID() == to.ID() {
				continue
			}
			if hasPath(from, to) {
				break
			}
		}
	}
}

func getReferences(to *node, allNodes rbtree.RbTree) []*node {
	if to.project.Project.ProjectReferences == nil {
		return []*node{}
	}

	dir := filepath.Dir(to.project.Path())

	result := make([]*node, len(to.project.Project.ProjectReferences))
	i := 0
	for _, ref := range to.project.Project.ProjectReferences {
		p := filepath.Join(dir, ref.Path())
		n := &node{fullPath: &p}
		from, ok := allNodes.Search(n)
		if ok {
			result[i] = from.(*node)
			i++
		}
	}
	return result[:i]
}
