package info

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"solt/internal/out"
	"solt/internal/ux"
	"solt/msvc"
	"solt/solution"
	"strconv"
	"strings"
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

func (d *display) solution(sl *msvc.VisualStudioSolution) {
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

func (d *display) showSectionsInfo(sections []*solution.Section) {
	var configurations = make(c9s.StringHashSet)
	var platforms = make(c9s.StringHashSet)

	for _, s := range sections {
		if s.Name != "SolutionConfigurationPlatforms" {
			continue
		}
		for _, item := range s.Items {
			parts := strings.Split(item.Key, "|")
			configuration := parts[0]
			platform := parts[1]
			configurations.Add(configuration)
			platforms.Add(platform)
		}
	}

	prn := newPrinter(d.margin, d.w)

	prn.print(configurations, "Configuration")

	d.p.Println()

	prn.print(platforms, "Platform")

	d.p.Println()
}
