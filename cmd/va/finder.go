package va

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"gonum.org/v1/gonum/graph/path"
)

type finder struct {
	g        *graph
	allPaths *path.AllShortest
}

func newFinder(g *graph) *finder {
	allPaths := g.allPaths()
	return &finder{allPaths: allPaths, g: g}
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

func (fi *finder) findAll() map[string]c9s.StringHashSet {
	result := make(map[string]c9s.StringHashSet)
	fi.g.foreach(func(n *node) {
		found, ok := fi.find(n.refs)

		if ok {
			result[n.String()] = found
		}
	})

	return result
}
