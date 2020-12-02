package cmd

import (
	"gonum.org/v1/gonum/graph"
	"solt/msvc"
)

type projectNode struct {
	id      int64
	project *msvc.MsbuildProject
}

func newProjectNode(id int64, prj *msvc.MsbuildProject) graph.Node {
	n := projectNode{
		id:      id,
		project: prj,
	}
	return &n
}

func (n *projectNode) ID() int64 {
	return n.id
}
