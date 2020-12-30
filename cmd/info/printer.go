package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
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
	format := p.m.Margin("%v\n")

	underline := api.NewCustomMarginer(len(name), "-").Margin("")
	p.p.Tprint(format, name)
	p.p.Tprint(format, underline)

	items := set.Items()
	sortfold.Strings(items)

	for _, k := range items {
		p.p.Tprint(format, k)
	}
	p.p.Flush()
}
