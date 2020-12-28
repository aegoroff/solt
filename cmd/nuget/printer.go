package nuget

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/api"
	"sort"
	"strings"
)

// pack defines nuget package descriptor
type pack struct {
	pkg      string
	versions c9s.StringHashSet
}

func newPack(id string, versions ...string) *pack {
	vs := make(c9s.StringHashSet)
	vs.AddRange(versions...)
	return &pack{
		pkg:      id,
		versions: vs,
	}
}

func copyPack(e *pack) *pack {
	return newPack(e.pkg, e.versions.Items()...)
}

func newNugetPrinter(p api.Printer) *nugetprint {
	np := nugetprint{
		p: p,
	}
	return &np
}

type nugetprint struct {
	p api.Printer
}

func (n *nugetprint) print(parent string, packs []*pack) {
	n.p.Cprint("\n <gray>%s</>\n", parent)

	const format = "  %v\t%v\n"
	n.p.Tprint(format, "Package", "Version")
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
