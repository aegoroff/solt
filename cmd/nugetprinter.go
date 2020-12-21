package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/gookit/color"
	"sort"
	"strings"
)

// pack defines nuget package descriptor
type pack struct {
	pkg      string
	versions c9s.StringHashSet
}

func newNugetPrinter(p printer) *nugetprint {
	np := nugetprint{
		p: p,
	}
	return &np
}

type nugetprint struct {
	p printer
}

func (n *nugetprint) print(parent string, packs []*pack) {
	n.p.setColor(color.FgGray)
	n.p.cprint("\n %s\n", parent)
	n.p.resetColor()

	const format = "  %v\t%v\n"
	n.p.tprint(format, "Package", "Version")
	n.p.tprint(format, "-------", "-------")

	sort.Slice(packs, func(i, j int) bool {
		return sortfold.CompareFold(packs[i].pkg, packs[j].pkg) < 0
	})

	for _, item := range packs {
		versions := item.versions.Items()
		sortfold.Strings(versions)
		n.p.tprint(format, item.pkg, strings.Join(versions, ", "))
	}
	n.p.flush()
}
