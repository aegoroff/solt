package graph

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/msvc"
	"strings"
)

// Node represents *msvc.MsbuildProject node in graph
type Node struct {
	id       int64
	project  *msvc.MsbuildProject
	fullPath *string
}

func newNode(id int64, prj *msvc.MsbuildProject) *Node {
	path := prj.Path()
	n := Node{
		id:       id,
		project:  prj,
		fullPath: &path,
	}
	return &n
}

// ID gets node's identifier
func (n *Node) ID() int64 {
	return n.id
}

// String gets project's full path
func (n *Node) String() string {
	return *n.fullPath
}

func (n *Node) Less(y rbtree.Comparable) bool {
	return sortfold.CompareFold(n.String(), y.(*Node).String()) < 0
}

func (n *Node) Equal(y rbtree.Comparable) bool {
	return strings.EqualFold(n.String(), y.(*Node).String())
}
