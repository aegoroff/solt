package info

import (
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
	tbl.AddHead("Projects", humanize.Comma(int64(t.result.projects)))

	types := make([]string, len(t.result.projectTypes))
	i := 0
	for k := range t.result.projectTypes {
		types[i] = k
		i++
	}

	sortfold.Strings(types)

	for _, name := range types {
		count := t.result.projectTypes[name]
		tbl.AddLine(name, humanize.Comma(int64(count)))
	}
	tbl.Print()
}
