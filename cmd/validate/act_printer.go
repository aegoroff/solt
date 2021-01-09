package validate

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/ux"
	"solt/internal/out"
)

type printer struct {
	prn out.Printer
}

func newPrinter(p out.Printer) actioner {
	return &printer{
		prn: p,
	}
}

func (v *printer) action(path string, refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}
	m1 := ux.NewMarginer(1)
	m3 := ux.NewMarginer(3)
	m5 := ux.NewMarginer(5)

	v.prn.Println()
	v.prn.Cprint(m1.Margin("Solution: <green>%s</>\n"), path)

	keys := make([]string, len(refs))
	i := 0
	for k := range refs {
		keys[i] = k
		i++
	}

	sortfold.Strings(keys)

	for _, k := range keys {
		v.prn.Cprint(m3.Margin("project <yellow>%s</> has redundant references:\n"), k)
		rrs := refs[k]

		items := rrs.SortedItems(sortfold.Strings)
		for _, s := range items {
			v.prn.Cprint(m5.Margin("<gray>%s</>\n"), s)
		}
		v.prn.Println()
	}
}
