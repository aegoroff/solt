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
	v.prn.Cprint(" Solution: <green>%s</>\n", name)

	projects := make([]string, 0, len(refs))
	for s := range refs {
		projects = append(projects, s)
	}

	sortfold.Strings(projects)

	hm := api.NewMarginer(3)
	rm := api.NewMarginer(5)
	for _, project := range projects {
		v.prn.Cprint(hm.Margin("project: <bold>%s</> has redundant references\n"), project)
		rrs := refs[project]

		items := rrs.Items()
		sortfold.Strings(items)
		for _, s := range items {
			v.prn.Cprint(rm.Margin("<gray>%s</>\n"), s)
		}
		v.prn.Cprint("\n")
	}
}
