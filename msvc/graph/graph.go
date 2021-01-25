package graph

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
	"path/filepath"
	"solt/msvc"
)

// Graph provides *msvc.MsbuildProject graph
type Graph struct {
	allNodes rbtree.RbTree
	g        *simple.DirectedGraph
	nextID   int64
}

// New creates new graph
func New(it *msvc.ProjectIterator) *Graph {
	gr := &Graph{
		g:        simple.NewDirectedGraph(),
		allNodes: rbtree.New(),
		nextID:   1,
	}

	it.Foreach(gr.newNode)
	ait := rbtree.NewWalkInorder(gr.allNodes)
	ait.Foreach(gr.newEdges)

	return gr
}

// AllPaths returns a shortest-path tree for shortest paths in the graph g.
// If the graph does not implement graph.Weighter, UniformCost is used.
// AllPaths will panic if g has a negative edge weight.
//
// The time complexity of AllPaths is O(|V|.|E|+|V|^2.log|V|).
func (gr *Graph) AllPaths() *path.AllShortest {
	paths := path.DijkstraAllPaths(gr.g)
	return &paths
}

// To returns all nodes in g that can reach directly to n.
func (gr *Graph) To(n *Node) []*Node {
	refs := gr.g.To(n.ID())
	nodes := make([]*Node, refs.Len())
	i := 0
	for refs.Next() {
		nodes[i] = refs.Node().(*Node)
		i++
	}
	return nodes
}

// Foreach enumerates all nodes in graph and call nodeFn on each node
func (gr *Graph) Foreach(nodeFn func(n *Node)) {
	it := rbtree.NewWalkInorder(gr.allNodes)
	it.Foreach(func(cmp rbtree.Comparable) {
		nodeFn(cmp.(*Node))
	})
}

func (gr *Graph) newNode(msbuild *msvc.MsbuildProject) {
	n := newNode(gr.nextID, msbuild)
	gr.allNodes.Insert(n)
	gr.nextID++
	gr.g.AddNode(n)
}

func (gr *Graph) newEdges(cmp rbtree.Comparable) {
	to := cmp.(*Node)
	if to.project.Project.ProjectReferences == nil {
		return
	}

	dir := filepath.Dir(to.project.Path())

	for _, ref := range to.project.Project.ProjectReferences {
		p := filepath.Join(dir, ref.Path())
		n := &Node{fullPath: &p}
		from, ok := gr.allNodes.Search(n)
		if ok {
			e := gr.g.NewEdge(from.(*Node), to)
			gr.g.SetEdge(e)
		}
	}
}
