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
		margin: margin,
	}
	return &np
}

type nugetprint struct {
	p      api.Printer
	column string
	margin int
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
	n.p.Cprint(n.newMargin("<gray>%s</>\n"), parent)

	format := n.newMargin("%v\t%v\n")
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

func (n *nugetprint) newMargin(s string) string {
	sb := strings.Builder{}
	for i := 0; i < n.margin; i++ {
		sb.WriteString(" ")
	}
	sb.WriteString(s)

	return sb.String()
}
