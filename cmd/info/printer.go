package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/api"
)

type printer struct {
	margin int
	w      api.Writable
}

func newPrinter(margin int, w api.Writable) *printer {
	return &printer{margin: margin, w: w}
}

func (p *printer) print(set c9s.StringHashSet, name string) {
	tbl := api.NewTabler(p.w, p.margin)
	tbl.AddHead(name)

	items := set.Items()
	sortfold.Strings(items)

	for _, k := range items {
		tbl.AddLine(k)
	}
	tbl.Print()
}
