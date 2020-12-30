package validate

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/msvc"
)

type node struct {
	id       int64
	project  *msvc.MsbuildProject
	fullPath *string
}

func newNode(id int64, prj *msvc.MsbuildProject) *node {
	n := node{
		id:       id,
		project:  prj,
		fullPath: &prj.Path,
	}
	return &n
}

func (n *node) ID() int64 {
	return n.id
}

func (n *node) String() string {
	return *n.fullPath
}

func (n *node) LessThan(y rbtree.Comparable) bool {
	return n.compare(y) < 0
}

func (n *node) EqualTo(y rbtree.Comparable) bool {
	return n.compare(y) == 0
}

func (n *node) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(n.String(), y.(*node).String())
}
