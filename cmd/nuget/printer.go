package nuget

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"solt/cmd/api"
	"sort"
	"strings"
)

func newNugetPrinter(p api.Printer) *nugetprint {
	np := nugetprint{
		p: p,
	}
	return &np
}

type nugetprint struct {
	p api.Printer
}

func (n *nugetprint) printTree(tree rbtree.RbTree, col string, head func(nf *nugetFolder) string) {
	it := rbtree.NewAscend(tree)

	it.Foreach(func(c rbtree.Comparable) {
		f := c.(*nugetFolder)
		n.print(head(f), col, f.packs)
	})
}

func (n *nugetprint) print(parent string, col string, packs []*pack) {
	n.p.Cprint("\n <gray>%s</>\n", parent)

	const format = "  %v\t%v\n"
	n.p.Tprint(format, col, "Version")
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
