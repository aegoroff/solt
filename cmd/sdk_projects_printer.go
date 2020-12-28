package cmd

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

func (v *sdkProjectsPrinter) action(sol string, refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}
	v.prn.Cprint(" Solution: <green>%s</>\n", sol)

	projects := make([]string, 0, len(refs))
	for s := range refs {
		projects = append(projects, s)
	}

	sortfold.Strings(projects)

	for _, project := range projects {
		v.prn.Cprint("   project: <bold>%s</> has redundant references\n", project)
		rrs := refs[project]

		items := rrs.Items()
		sortfold.Strings(items)
		for _, s := range items {
			v.prn.Cprint("     <gray>%s</>\n", s)
		}
	}
}
