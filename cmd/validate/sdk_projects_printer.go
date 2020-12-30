package validate

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
	"solt/cmd/api"
)

type sdkProjectsPrinter struct {
	prn api.Printer
}

func newSdkProjectsPrinter(p api.Printer) sdkActioner {
	return &sdkProjectsPrinter{
		prn: p,
	}
}

func (v *sdkProjectsPrinter) action(name string, refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}
	sm := api.NewMarginer(1)
	pm := api.NewMarginer(3)
	rm := api.NewMarginer(5)

	v.prn.Cprint("\n")
	v.prn.Cprint(sm.Margin("Solution: <green>%s</>\n"), name)

	projects := make([]string, len(refs))
	i := 0
	for s := range refs {
		projects[i] = s
		i++
	}

	sortfold.Strings(projects)

	for _, project := range projects {
		v.prn.Cprint(pm.Margin("project: <bold>%s</> has redundant references\n"), project)
		rrs := refs[project]

		items := rrs.Items()
		sortfold.Strings(items)
		for _, s := range items {
			v.prn.Cprint(rm.Margin("<gray>%s</>\n"), s)
		}
		v.prn.Cprint("\n")
	}
}
