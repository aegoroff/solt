package va

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

type graph struct {
	allNodes rbtree.RbTree
	g        *simple.DirectedGraph
}

func (gr *graph) allPaths() *path.AllShortest {
	paths := path.DijkstraAllPaths(gr.g)
	return &paths
}

func (gr *graph) foreach(f func(n *node)) {
	gn := gr.g.Nodes()

	for gn.Next() {
		n := gn.Node().(*node)
		f(n)
	}
}

func newGraph(sln *msvc.VisualStudioSolution, projects rbtree.RbTree) *graph {
	g := simple.NewDirectedGraph()
	allNodes := rbtree.New()

	gr := &graph{
		g:        g,
		allNodes: allNodes,
	}

	gr.newNodes(sln, projects)

	gr.newEdges()

	return gr
}

func (gr *graph) newNodes(sln *msvc.VisualStudioSolution, projects rbtree.RbTree) {
	solutionPath := filepath.Dir(sln.Path())
	ix := int64(1)
	for _, prj := range sln.Solution.Projects {
		if prj.TypeID == solution.IDSolutionFolder {
			continue
		}

		p := msvc.NewMsbuildProject(filepath.Join(solutionPath, prj.Path))

		msbuild, ok := projects.Search(p)
		if !ok {
			continue
		}

		n := newNode(ix, msbuild.(*msvc.MsbuildProject))
		gr.allNodes.Insert(n)
		ix++
		gr.g.AddNode(n)
	}
}

func (gr *graph) newEdges() {
	gn := gr.g.Nodes()

	for gn.Next() {
		to := gn.Node().(*node)
		to.refs = gr.references(to)
		for _, from := range to.refs {
			e := gr.g.NewEdge(from, to)
			gr.g.SetEdge(e)
		}
	}
}

func (gr *graph) references(to *node) []*node {
	if to.project.Project.ProjectReferences == nil {
		return []*node{}
	}

	dir := filepath.Dir(to.project.Path())

	result := make([]*node, len(to.project.Project.ProjectReferences))
	i := 0
	for _, ref := range to.project.Project.ProjectReferences {
		p := filepath.Join(dir, ref.Path())
		n := &node{fullPath: &p}
		from, ok := gr.allNodes.Search(n)
		if ok {
			result[i] = from.(*node)
			i++
		}
	}
	return result[:i]
}
