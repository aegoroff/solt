package info

import (
	"solt/internal/out"
	"solt/internal/ux"
	"solt/msvc"
	"solt/solution"
	"strconv"
)

type display struct {
	p      out.Printer
	margin int
	w      out.Writable
}

func newDisplay(p out.Printer, w out.Writable) *display {
	return &display{
		p:      p,
		margin: 2,
		w:      w,
	}
}

func (d *display) Solution(sl *msvc.VisualStudioSolution) {
	sln := sl.Solution

	d.p.Cprint(" <gray>%s</>\n", sl.Path())

	tbl := ux.NewTabler(d.w, d.margin)

	tbl.AddLine("Header", sln.Header)
	tbl.AddLine("Product", sln.Comment)
	tbl.AddLine("Visual Studio Version", sln.VisualStudioVersion)
	tbl.AddLine("Minimum Visual Studio Version", sln.MinimumVisualStudioVersion)

	tbl.Print()

	d.p.Println()

	d.showProjectsInfo(sln.Projects)
	d.showSectionsInfo(sln.GlobalSections)
}

func (d *display) showProjectsInfo(projects []*solution.Project) {
	var byType = make(map[string]int)

	for _, p := range projects {
		byType[p.Type]++
	}

	tbl := ux.NewTabler(d.w, d.margin)
	tbl.AddHead("Project type", "Count")

	for k, v := range byType {
		tbl.AddLine(k, strconv.Itoa(v))
	}
	tbl.Print()
	d.p.Println()
}

func (d *display) showSectionsInfo(sections sections) {
	confPlat := newConfigurationPlatform()

	sections.foreach(confPlat)

	prn := newPrinter(d.margin, d.w)

	prn.print(confPlat.configurations, "Configuration")

	d.p.Println()

	prn.print(confPlat.platforms, "Platform")

	d.p.Println()
}
