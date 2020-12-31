package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/cmd/api"
	"sort"
	"strings"
)

func newNugetPrinter(p api.Printer, w api.Writable, column string, margin int) *nugetprint {
	np := nugetprint{
		p:      p,
		w:      w,
		column: column,
		m:      api.NewMarginer(margin),
	}
	return &np
}

type nugetprint struct {
	p      api.Printer
	w      api.Writable
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

	tbl := api.NewTabler(n.w, n.m.Value())
	tbl.AddHead(n.column, "Version")

	sort.Slice(packs, func(i, j int) bool {
		return sortfold.CompareFold(packs[i].pkg, packs[j].pkg) < 0
	})

	for _, item := range packs {
		versions := item.versions.Items()
		sortfold.Strings(versions)

		tbl.AddLine(item.pkg, strings.Join(versions, ", "))
	}
	tbl.Print()
}
