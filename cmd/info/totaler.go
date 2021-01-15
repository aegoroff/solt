package info

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"solt/internal/out"
	"solt/internal/ux"
	"solt/msvc"
	"sort"
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
	p.Cprint(" <red>Totals:</>\n")

	tbl := ux.NewTabler(w, 2)
	tbl.AddLine("Solutions", humanize.Comma(int64(t.result.solutions)))
	tbl.AddLine("Projects", humanize.Comma(int64(t.result.projects)))
	tbl.AddLine("", "")

	type nv struct {
		name string
		val  int
	}

	types := make([]nv, 0, len(t.result.projectTypes))
	for k, v := range t.result.projectTypes {
		types = append(types, nv{name: k, val: v})
	}

	sort.Slice(types, func(i, j int) bool {
		return types[i].val > types[j].val
	})

	tbl.AddHead("Project type", "Count", "Percent")
	for _, tt := range types {
		countS := humanize.Comma(int64(tt.val))
		percent := t.percent(tt.val)
		percentS := fmt.Sprintf("%.2f%%", percent)
		tbl.AddLine(tt.name, countS, percentS)
	}
	tbl.Print()
}

func (t *totaler) percent(value int) float64 {
	if t.result.projects == 0 {
		return 0
	}
	return (float64(value) / float64(t.result.projects)) * 100
}
