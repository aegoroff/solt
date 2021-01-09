package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"strings"
)

type folder struct {
	path    string
	sources []string
	packs   []*pack
}

func (n *folder) Less(y rbtree.Comparable) bool {
	return sortfold.CompareFold(n.path, y.(*folder).path) < 0
}

func (n *folder) Equal(y rbtree.Comparable) bool {
	return strings.EqualFold(n.path, y.(*folder).path)
}

func newNugetFolder(p string, packs []*pack, src []string) rbtree.Comparable {
	nf := folder{
		path:    p,
		packs:   packs,
		sources: src,
	}
	return &nf
}
