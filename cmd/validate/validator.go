package validate

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/cmd/fw"
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
	g, allNodes := va.newSolutionGraph(sol)

	redundants := va.findRedundants(g, allNodes)
	if len(redundants) > 0 {
		va.tt.problemSolutions++
		va.tt.problemProjects += int64(len(redundants))
		for _, set := range redundants {
			va.tt.redundantRefs += int64(set.Count())
		}
	}

	va.act.action(sol.Path(), redundants)
}

func (va *validator) newSolutionGraph(sln *msvc.VisualStudioSolution) (*simple.DirectedGraph, rbtree.RbTree) {
	g, nodes := va.createGraphNodes(sln)
	va.createGraphEdges(g, nodes)
	return g, nodes
}

func (va *validator) createGraphNodes(sln *msvc.VisualStudioSolution) (*simple.DirectedGraph, rbtree.RbTree) {
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

	return g, nodes
}

func (va *validator) createGraphEdges(g *simple.DirectedGraph, nodes rbtree.RbTree) {
	gn := g.Nodes()

	for gn.Next() {
		to := gn.Node().(*node)
		refs := va.getReferences(to, nodes)
		for _, ref := range refs {
			e := g.NewEdge(ref, to)
			g.SetEdge(e)
		}
	}
}

func (va *validator) findRedundants(g *simple.DirectedGraph, allNodes rbtree.RbTree) map[string]c9s.StringHashSet {
	allPaths := path.DijkstraAllPaths(g)
	result := make(map[string]c9s.StringHashSet)

	gn := g.Nodes()

	for gn.Next() {
		project := gn.Node().(*node)

		refs := va.getReferences(project, allNodes)

		rrs := make(c9s.StringHashSet)

		allPairs(refs, func(from *node, to *node) {
			paths, _ := allPaths.AllBetween(from.ID(), to.ID())
			if len(paths) > 0 {
				rrs.Add(from.String())
			}
		})

		if rrs.Count() > 0 {
			result[project.String()] = rrs
		}
	}

	return result
}

func allPairs(nodes []*node, action func(*node, *node)) {
	for _, from := range nodes {
		for _, to := range nodes {
			if from.ID() == to.ID() {
				continue
			}
			action(from, to)
		}
	}
}

func (*validator) getReferences(to *node, allNodes rbtree.RbTree) []*node {
	if to.project.Project.ProjectReferences == nil {
		return []*node{}
	}

	dir := filepath.Dir(to.project.Path())

	var result []*node
	for _, ref := range to.project.Project.ProjectReferences {
		p := filepath.Join(dir, ref.Path())
		n := &node{fullPath: &p}
		from, ok := allNodes.Search(n)
		if ok {
			result = append(result, from.(*node))
		}
	}
	return result
}
