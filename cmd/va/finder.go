package va

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

type finder struct {
	allPaths *path.AllShortest
}

func newFinder(g *simple.DirectedGraph) *finder {
	allPaths := path.DijkstraAllPaths(g)
	return &finder{allPaths: &allPaths}
}

func (fi *finder) hasPath(from *node, to *node) bool {
	if from.ID() == to.ID() {
		return false
	}
	paths, _ := fi.allPaths.AllBetween(from.ID(), to.ID())
	return len(paths) > 0
}

func (fi *finder) find(nodes []*node) (c9s.StringHashSet, bool) {
	found := c9s.NewStringHashSet()
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
