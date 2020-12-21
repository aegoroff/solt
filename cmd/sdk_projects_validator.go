package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

type sdkProjectsValidator struct {
	prn         printer
	fs          afero.Fs
	sourcesPath string
	actioner    sdkActioner
}

func newSdkProjectsValidator(fs afero.Fs, p printer, sourcesPath string, h sdkActioner) *sdkProjectsValidator {
	return &sdkProjectsValidator{
		prn:         p,
		fs:          fs,
		sourcesPath: sourcesPath,
		actioner:    h,
	}
}

func (m *sdkProjectsValidator) validate() {
	foldersTree := msvc.ReadSolutionDir(m.sourcesPath, m.fs)

	solutions, allProjects := msvc.SelectSolutionsAndProjects(foldersTree)

	sdkProjects := m.onlySdkProjects(allProjects)

	for _, sol := range solutions {
		g, allNodes := m.newSolutionGraph(sol, sdkProjects)

		redundants := m.findRedundants(g, allNodes)

		m.actioner.action(sol.Path, redundants)
	}
}

func (*sdkProjectsValidator) onlySdkProjects(allProjects []*msvc.MsbuildProject) rbtree.RbTree {
	tree := rbtree.NewRbTree()

	for _, project := range allProjects {
		if !project.Project.IsSdkProject() {
			continue
		}
		tree.Insert(project)
	}
	return tree
}

func (m *sdkProjectsValidator) newSolutionGraph(sln *msvc.VisualStudioSolution, sdkTree rbtree.RbTree) (*simple.DirectedGraph, rbtree.RbTree) {
	g, nodes := m.createGraphNodes(sln, sdkTree)
	m.createGraphEdges(g, nodes)
	return g, nodes
}

func (*sdkProjectsValidator) createGraphNodes(sln *msvc.VisualStudioSolution, sdkTree rbtree.RbTree) (*simple.DirectedGraph, rbtree.RbTree) {
	solutionPath := filepath.Dir(sln.Path)
	g := simple.NewDirectedGraph()
	nodes := rbtree.NewRbTree()
	ix := int64(1)
	for _, prj := range sln.Solution.Projects {
		if prj.TypeID == solution.IDSolutionFolder {
			continue
		}

		p := &msvc.MsbuildProject{
			Path: filepath.Join(solutionPath, prj.Path),
		}

		msbuild, ok := sdkTree.Search(p)
		if !ok {
			continue
		}

		n := newProjectNode(ix, msbuild.Key().(*msvc.MsbuildProject))
		nodes.Insert(n)
		ix++
		g.AddNode(n)
	}

	return g, nodes
}

func (m *sdkProjectsValidator) createGraphEdges(g *simple.DirectedGraph, nodes rbtree.RbTree) {
	gn := g.Nodes()

	for gn.Next() {
		to := gn.Node().(*projectNode)
		refs := m.getReferences(to, nodes)
		for _, ref := range refs {
			e := g.NewEdge(ref, to)
			g.SetEdge(e)
		}
	}
}

func (m *sdkProjectsValidator) findRedundants(g *simple.DirectedGraph, allNodes rbtree.RbTree) map[string]c9s.StringHashSet {
	allPaths := path.DijkstraAllPaths(g)
	result := make(map[string]c9s.StringHashSet)

	gn := g.Nodes()

	for gn.Next() {
		project := gn.Node().(*projectNode)

		refs := m.getReferences(project, allNodes)

		rrs := make(c9s.StringHashSet)

		allPairs(refs, func(from *projectNode, to *projectNode) {
			paths, _ := allPaths.AllBetween(from.ID(), to.ID())
			if paths != nil && len(paths) > 0 {
				rrs.Add(from.String())
			}
		})

		if rrs.Count() > 0 {
			result[project.String()] = rrs
		}
	}

	return result
}

func allPairs(nodes []*projectNode, action func(*projectNode, *projectNode)) {
	for _, from := range nodes {
		for _, to := range nodes {
			if from.ID() == to.ID() {
				continue
			}
			action(from, to)
		}
	}
}

func (*sdkProjectsValidator) getReferences(to *projectNode, allNodes rbtree.RbTree) []*projectNode {
	if to.project.Project.ProjectReferences == nil {
		return []*projectNode{}
	}

	dir := filepath.Dir(to.project.Path)

	var result []*projectNode
	for _, ref := range to.project.Project.ProjectReferences {
		p := filepath.Join(dir, ref.Path)
		n := &projectNode{fullPath: &p}
		from, ok := allNodes.Search(n)
		if ok {
			result = append(result, from.Key().(*projectNode))
		}
	}
	return result
}
