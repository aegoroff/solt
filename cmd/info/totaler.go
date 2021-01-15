package info

import (
	"fmt"
	"github.com/akutz/sortfold"
	"github.com/dustin/go-humanize"
	"solt/internal/out"
	"solt/internal/ux"
	"solt/msvc"
)

type totaler struct {
	result *totals
	grp    *projectGroupper
}

func newTotaler(grp *projectGroupper) *totaler {
	return &totaler{
		result: &totals{projectTypes: make(map[string]int)},
		grp:    grp,
	}
}

func (t *totaler) Solution(*msvc.VisualStudioSolution) {
	t.result.solutions++
	for k, v := range t.groupped() {
		t.result.projectTypes[k] += v
		t.result.projects += v
	}
}

func (t *totaler) groupped() map[string]int {
	return t.grp.ByType()
}

func (t *totaler) display(p out.Printer, w out.Writable) {
	p.Cprint(" <blue>Totals:</>\n")

	tbl := ux.NewTabler(w, 2)
	tbl.AddLine("Solutions", humanize.Comma(int64(t.result.solutions)))
	tbl.AddLine("Projects", humanize.Comma(int64(t.result.projects)))
	tbl.AddLine("", "")

	types := make([]string, len(t.result.projectTypes))
	i := 0
	for k := range t.result.projectTypes {
		types[i] = k
		i++
	}

	sortfold.Strings(types)

	tbl.AddHead("Project type", "Count", "Percent")
	for _, name := range types {
		count := t.result.projectTypes[name]
		countS := humanize.Comma(int64(count))
		percent := t.percent(count)
		percentS := fmt.Sprintf("%.2f%%", percent)
		tbl.AddLine(name, countS, percentS)
	}
	tbl.Print()
}

func (t *totaler) percent(value int) float64 {
	if t.result.projects == 0 {
		return 0
	}
	return (float64(value) / float64(t.result.projects)) * 100
}
