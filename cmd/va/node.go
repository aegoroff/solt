package va

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/msvc"
	"strings"
)

type node struct {
	id       int64
	project  *msvc.MsbuildProject
	fullPath *string
}

func newNode(id int64, prj *msvc.MsbuildProject) *node {
	path := prj.Path()
	n := node{
		id:       id,
		project:  prj,
		fullPath: &path,
	}
	return &n
}

func (n *node) ID() int64 {
	return n.id
}

func (n *node) String() string {
	return *n.fullPath
}

func (n *node) Less(y rbtree.Comparable) bool {
	return sortfold.CompareFold(n.String(), y.(*node).String()) < 0
}

func (n *node) Equal(y rbtree.Comparable) bool {
	return strings.EqualFold(n.String(), y.(*node).String())
}
