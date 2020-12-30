package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/api"
)

type printer struct {
	margin int
	p      api.Printer
}

func newPrinter(margin int, p api.Printer) *printer {
	return &printer{margin: margin, p: p}
}

func (p *printer) print(set c9s.StringHashSet, name string) {
	tbl := api.NewTabler(p.p, p.margin)
	tbl.AddHead(name)

	items := set.Items()
	sortfold.Strings(items)

	for _, k := range items {
		tbl.AddLine(k)
	}
	tbl.Print()
}
