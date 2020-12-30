package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"github.com/cheynewallace/tabby"
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

	t := tabby.NewCustom(n.p.Twriter())

	const ver = "Version"
	cunder := api.NewUnderline(n.column)
	vunder := api.NewUnderline(ver)

	t.AddLine(n.m.Margin(n.column), ver)
	t.AddLine(n.m.Margin(cunder), vunder)

	sort.Slice(packs, func(i, j int) bool {
		return sortfold.CompareFold(packs[i].pkg, packs[j].pkg) < 0
	})

	for _, item := range packs {
		versions := item.versions.Items()
		sortfold.Strings(versions)

		t.AddLine(n.m.Margin(item.pkg), strings.Join(versions, ", "))
	}

	t.Print()
}
