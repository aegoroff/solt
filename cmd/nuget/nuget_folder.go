package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"strings"
)

type nugetFolder struct {
	path    string
	sources []string
	packs   []*pack
}

func (n *nugetFolder) Less(y rbtree.Comparable) bool {
	return sortfold.CompareFold(n.path, y.(*nugetFolder).path) < 0
}

func (n *nugetFolder) Equal(y rbtree.Comparable) bool {
	return strings.EqualFold(n.path, y.(*nugetFolder).path)
}

func newNugetFolder(p string, packs []*pack, src []string) rbtree.Comparable {
	nf := nugetFolder{
		path:    p,
		packs:   packs,
		sources: src,
	}
	return &nf
}
