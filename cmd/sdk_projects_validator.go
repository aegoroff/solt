package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/akutz/sortfold"
)

type sdkProjectsValidator struct {
	prn printer
}

func newsdkProjectsValidator(p printer) sdkModuleHandler {
	return &sdkProjectsValidator{
		prn: p,
	}
}

func (v *sdkProjectsValidator) handle(sol string, refs map[string]c9s.StringHashSet) {
	if len(refs) == 0 {
		return
	}
	v.prn.cprint(" Solution: <green>%s</>\n", sol)

	projects := make([]string, 0, len(refs))
	for s := range refs {
		projects = append(projects, s)
	}

	sortfold.Strings(projects)

	for _, project := range projects {
		v.prn.cprint("   project: <bold>%s</> has redundant references\n", project)
		rrs := refs[project]

		items := rrs.Items()
		sortfold.Strings(items)
		for _, s := range items {
			v.prn.cprint("     <gray>%s</>\n", s)
		}
	}
}