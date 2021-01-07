package validate

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/fw"
)

type printer struct {
	prn fw.Printer
}

func newPrinter(p fw.Printer) actioner {
	return &printer{
		prn: p,
	}
}

func (v *printer) action(path string, refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}
	sm := fw.NewMarginer(1)
	pm := fw.NewMarginer(3)
	rm := fw.NewMarginer(5)

	v.prn.Println()
	v.prn.Cprint(sm.Margin("Solution: <green>%s</>\n"), path)

	projects := make([]string, len(refs))
	i := 0
	for s := range refs {
		projects[i] = s
		i++
	}

	sortfold.Strings(projects)

	for _, project := range projects {
		v.prn.Cprint(pm.Margin("project <yellow>%s</> has redundant references:\n"), project)
		rrs := refs[project]

		items := rrs.SortedItems(sortfold.Strings)
		for _, s := range items {
			v.prn.Cprint(rm.Margin("<gray>%s</>\n"), s)
		}
		v.prn.Println()
	}
}
