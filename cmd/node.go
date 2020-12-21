package cmd

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/msvc"
)

type projectNode struct {
	id       int64
	project  *msvc.MsbuildProject
	fullPath *string
}

func newProjectNode(id int64, prj *msvc.MsbuildProject) *projectNode {
	n := projectNode{
		id:       id,
		project:  prj,
		fullPath: &prj.Path,
	}
	return &n
}

func (n *projectNode) ID() int64 {
	return n.id
}

func (n *projectNode) String() string {
	return *n.fullPath
}

func (n *projectNode) LessThan(y rbtree.Comparable) bool {
	return n.compare(y) < 0
}

func (n *projectNode) EqualTo(y rbtree.Comparable) bool {
	return n.compare(y) == 0
}

func (n *projectNode) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(n.String(), y.(*projectNode).String())
}
