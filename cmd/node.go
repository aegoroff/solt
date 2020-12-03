package cmd

import (
	"solt/msvc"
)

type projectNode struct {
	id      int64
	project *msvc.MsbuildProject
}

func newProjectNode(id int64, prj *msvc.MsbuildProject) *projectNode {
	n := projectNode{
		id:      id,
		project: prj,
	}
	return &n
}

func (n *projectNode) ID() int64 {
	return n.id
}

func (n *projectNode) String() string {
	return n.project.Path
}
