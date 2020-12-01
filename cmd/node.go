package cmd

import "gonum.org/v1/gonum/graph"

type projectNode struct {
	id   int64
	path string
}

func newProjectNode(id int64, path string) graph.Node {
	n := projectNode{
		id:   id,
		path: path,
	}
	return &n
}

func (n *projectNode) ID() int64 {
	return n.id
}
