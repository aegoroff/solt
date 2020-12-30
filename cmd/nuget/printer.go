package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/cmd/api"
	"sort"
	"strings"
)

func newNugetPrinter(p api.Printer, column string, margin int) *nugetprint {
	np := nugetprint{
		p:      p,
		column: column,
		m:      api.NewMarginer(margin),
	}
	return &np
}

type nugetprint struct {
	p      api.Printer
	column string
	m      *api.Marginer
}

func (n *nugetprint) printTree(tree rbtree.RbTree, head func(nf *folder) string) {
	it := rbtree.NewAscend(tree)

	it.Foreach(func(c rbtree.Comparable) {
		f := c.(*folder)
		n.print(head(f), f.packs)
	})
}

func (n *nugetprint) print(parent string, packs []*pack) {
	n.p.Cprint("\n")
	n.p.Cprint(n.m.Margin("<gray>%s</>\n"), parent)

	format := n.m.Margin("%v\t%v\n")
	n.p.Tprint(format, n.column, "Version")
	n.p.Tprint(format, "-------", "-------")

	sort.Slice(packs, func(i, j int) bool {
		return sortfold.CompareFold(packs[i].pkg, packs[j].pkg) < 0
	})

	for _, item := range packs {
		versions := item.versions.Items()
		sortfold.Strings(versions)
		n.p.Tprint(format, item.pkg, strings.Join(versions, ", "))
	}
	n.p.Flush()
}
