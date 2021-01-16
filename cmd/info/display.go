package info

import (
	"solt/internal/out"
	"solt/internal/ux"
	"solt/msvc"
	"sort"
	"strconv"
)

type display struct {
	p      out.Printer
	margin int
	w      out.Writable
	grp    *projectGroupper
}

func newDisplay(p out.Printer, w out.Writable, grp *projectGroupper) *display {
	return &display{
		p:      p,
		margin: 2,
		w:      w,
		grp:    grp,
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

	d.showProjectsInfo()
	d.showSectionsInfo(sln.GlobalSections)
}

func (d *display) showProjectsInfo() {
	tbl := ux.NewTabler(d.w, d.margin+1)
	tbl.AddHead("Project type", "Count")

	type kv struct {
		k string
		v int
	}
	var kvs []kv
	for k, v := range d.groupped() {
		kvs = append(kvs, kv{k, v})
	}
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].k < kvs[j].k
	})
	for _, v := range kvs {
		tbl.AddLine(v.k, strconv.Itoa(v.v))
	}
	tbl.Print()
	d.p.Println()
}

func (d *display) groupped() map[string]int {
	return d.grp.ByType()
}

func (d *display) showSectionsInfo(sections sections) {
	confPlat := newConfigurationPlatform()

	sections.foreach(confPlat)

	prn := newPrinter(d.margin+1, d.w)

	prn.print(confPlat.configurations, "Configuration")

	d.p.Println()

	prn.print(confPlat.platforms, "Platform")

	d.p.Println()
}
