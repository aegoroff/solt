package validate

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/fw"
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
	m1 := fw.NewMarginer(1)
	m3 := fw.NewMarginer(3)
	m5 := fw.NewMarginer(5)

	v.prn.Println()
	v.prn.Cprint(m1.Margin("Solution: <green>%s</>\n"), path)

	projects := make([]string, len(refs))
	i := 0
	for s := range refs {
		projects[i] = s
		i++
	}

	sortfold.Strings(projects)

	for _, project := range projects {
		v.prn.Cprint(m3.Margin("project <yellow>%s</> has redundant references:\n"), project)
		rrs := refs[project]

		items := rrs.SortedItems(sortfold.Strings)
		for _, s := range items {
			v.prn.Cprint(m5.Margin("<gray>%s</>\n"), s)
		}
		v.prn.Println()
	}
}
