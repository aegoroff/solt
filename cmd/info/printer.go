package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"github.com/cheynewallace/tabby"
	"solt/cmd/api"
)

type printer struct {
	m *api.Marginer
	p api.Printer
}

func newPrinter(m *api.Marginer, p api.Printer) *printer {
	return &printer{m: m, p: p}
}

func (p *printer) print(set c9s.StringHashSet, name string) {
	t := tabby.NewCustom(p.p.Twriter())

	underline := api.NewCustomMarginer(len(name), "-").Margin("")

	t.AddLine(p.m.Margin(name))
	t.AddLine(p.m.Margin(underline))

	items := set.Items()
	sortfold.Strings(items)

	for _, k := range items {
		t.AddLine(p.m.Margin(k))
	}
	t.Print()
}
