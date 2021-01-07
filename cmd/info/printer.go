package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/fw"
)

type printer struct {
	margin int
	w      fw.Writable
}

func newPrinter(margin int, w fw.Writable) *printer {
	return &printer{margin: margin, w: w}
}

func (p *printer) print(set c9s.StringHashSet, name string) {
	tbl := fw.NewTabler(p.w, p.margin)
	tbl.AddHead(name)

	items := set.SortedItems(sortfold.Strings)

	for _, k := range items {
		tbl.AddLine(k)
	}
	tbl.Print()
}
