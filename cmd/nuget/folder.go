package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
)

type nugetFolder struct {
	path    string
	sources []string
	packs   []*pack
}

func (n *nugetFolder) LessThan(y rbtree.Comparable) bool {
	return n.compare(y) < 0
}

func (n *nugetFolder) EqualTo(y rbtree.Comparable) bool {
	return n.compare(y) == 0
}

func (n *nugetFolder) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(n.path, y.(*nugetFolder).path)
}

func newNugetFolder(p string, packs []*pack, src []string) rbtree.Comparable {
	nf := nugetFolder{
		path:    p,
		packs:   packs,
		sources: src,
	}
	return &nf
}