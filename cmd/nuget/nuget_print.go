package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/internal/out"
	"solt/internal/ux"
	"sort"
	"strings"
)

func newNugetPrint(p out.Printer, w out.Writable, column string, margin int) *nugetprint {
	np := nugetprint{
		p:      p,
		w:      w,
		column: column,
		m:      ux.NewMarginer(margin),
	}
	return &np
}

type nugetprint struct {
	p      out.Printer
	w      out.Writable
	column string
	m      *ux.Marginer
}

func (n *nugetprint) printTree(tree rbtree.RbTree, head func(nf *nugetFolder) string) {
	it := rbtree.NewAscend(tree)

	it.Foreach(func(c rbtree.Comparable) {
		f := c.(*nugetFolder)
		n.print(head(f), f.packs)
	})
}

func (n *nugetprint) print(parent string, packs []*pack) {
	n.p.Cprint("\n")
	n.p.Cprint(n.m.Margin("<gray>%s</>\n"), parent)

	tbl := ux.NewTabler(n.w, n.m.Value()+1)
	tbl.AddHead(n.column, "Version")

	sort.Slice(packs, func(i, j int) bool {
		return sortfold.CompareFold(packs[i].pkg, packs[j].pkg) < 0
	})

	for _, item := range packs {
		versions := item.versions.SortedItems(sortfold.Strings)

		tbl.AddLine(item.pkg, strings.Join(versions, ", "))
	}
	tbl.Print()
}
