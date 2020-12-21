package cmd

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"path/filepath"
	"solt/msvc"
)

type projectNode struct {
	id       int64
	project  *msvc.MsbuildProject
	fullPath string
}

func newProjectNode(id int64, prj *msvc.MsbuildProject, basePath string) *projectNode {
	n := projectNode{
		id:       id,
		project:  prj,
		fullPath: filepath.Join(basePath, prj.Path),
	}
	return &n
}

func (n *projectNode) ID() int64 {
	return n.id
}

func (n *projectNode) String() string {
	return n.project.Path
}

func (n *projectNode) LessThan(y rbtree.Comparable) bool {
	return n.compare(y) < 0
}

func (n *projectNode) EqualTo(y rbtree.Comparable) bool {
	return n.compare(y) == 0
}

func (x *projectNode) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(x.fullPath, y.(*projectNode).fullPath)
}
