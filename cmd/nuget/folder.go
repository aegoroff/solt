package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
)

type folder struct {
	path    string
	sources []string
	packs   []*pack
}

func (n *folder) LessThan(y rbtree.Comparable) bool {
	return n.compare(y) < 0
}

func (n *folder) EqualTo(y rbtree.Comparable) bool {
	return n.compare(y) == 0
}

func (n *folder) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(n.path, y.(*folder).path)
}

func newNugetFolder(p string, packs []*pack, src []string) rbtree.Comparable {
	nf := folder{
		path:    p,
		packs:   packs,
		sources: src,
	}
	return &nf
}
