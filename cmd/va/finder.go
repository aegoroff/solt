package va

import (
	"solt/msvc/graph"

	c9s "github.com/aegoroff/godatastruct/collections"
	"gonum.org/v1/gonum/graph/path"
)

type finder struct {
	g        *graph.Graph
	allPaths *path.AllShortest
}

func newFinder(g *graph.Graph) *finder {
	allPaths := g.AllPaths()
	return &finder{allPaths: allPaths, g: g}
}

func (fi *finder) hasPath(from *graph.Node, to *graph.Node) bool {
	if from.ID() == to.ID() {
		return false
	}
	paths, _ := fi.allPaths.AllBetween(from.ID(), to.ID())
	return len(paths) > 0
}

func (fi *finder) find(n *graph.Node) (c9s.HashSet[string], bool) {
	nodes := fi.g.To(n)
	found := c9s.NewHashSet[string]()
	for _, from := range nodes {
		for _, to := range nodes {
			if fi.hasPath(from, to) {
				found.Add(from.String())
				break
			}
		}
	}
	return found, found.Count() > 0
}

func (fi *finder) findAll() map[string]c9s.HashSet[string] {
	result := make(map[string]c9s.HashSet[string])
	fi.g.Foreach(func(n *graph.Node) {
		found, ok := fi.find(n)

		if ok {
			result[n.String()] = found
		}
	})

	return result
}
