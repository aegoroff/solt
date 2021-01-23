package va

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
)

type graph struct {
	allNodes rbtree.RbTree
	g        *simple.DirectedGraph
	nextID   int64
}

func newGraph(sln *msvc.VisualStudioSolution, it *sdkIterator) *graph {
	gr := &graph{
		g:        simple.NewDirectedGraph(),
		allNodes: rbtree.New(),
		nextID:   1,
	}

	it.foreach(sln, gr.newNode)
	ait := rbtree.NewWalkInorder(gr.allNodes)
	ait.Foreach(gr.newEdge)

	return gr
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

func (gr *graph) newNode(msbuild *msvc.MsbuildProject) {
	n := newNode(gr.nextID, msbuild)
	gr.allNodes.Insert(n)
	gr.nextID++
	gr.g.AddNode(n)
}

func (gr *graph) newEdge(cmp rbtree.Comparable) {
	to := cmp.(*node)
	to.refs = gr.references(to)
	for _, from := range to.refs {
		e := gr.g.NewEdge(from, to)
		gr.g.SetEdge(e)
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
