package info

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"solt/cmd/fw"
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
		result: &totals{
			projectTypes: rbtree.New(),
		},
		grp: grp,
	}
}

func (t *totaler) Solution(*msvc.VisualStudioSolution) {
	t.result.solutions++
	for k, v := range t.groupped() {
		t.updateType(k, int64(v))
	}
}

func (t *totaler) updateType(k string, v int64) {
	t.result.projects += v
	key := newTypeStat(k)
	n, ok := t.result.projectTypes.Search(newTypeStat(k))
	if ok {
		ts := n.(*typeStat)
		ts.solutions++
		ts.count += v
	} else {
		key.solutions = 1
		key.count = v
		t.result.projectTypes.Insert(key)
	}
}

func (t *totaler) groupped() map[string]int {
	return t.grp.ByType()
}

func (t *totaler) display(p out.Printer, w out.Writable) {
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)
	tbl.AddLine("Solutions", humanize.Comma(t.result.solutions))
	tbl.AddLine("Projects", humanize.Comma(t.result.projects))
	tbl.AddLine("", "")

	const percentH = "%     "
	tbl.AddHead("Project type", "Count", percentH, "Solutions", percentH)

	it := rbtree.NewAscend(t.result.projectTypes)

	it.Foreach(func(n rbtree.Comparable) {
		ts := n.(*typeStat)
		percentS := fmt.Sprintf("%.2f%%", t.percentProjects(ts.count))
		solPercentS := fmt.Sprintf("%.2f%%", t.percentSolutions(ts.solutions))
		tbl.AddLine(ts.name, ts.Count(), percentS, ts.Solutions(), solPercentS)
	})

	tbl.Print()
}

func (t *totaler) percentProjects(value int64) float64 {
	return fw.Percent(value, t.result.projects)
}

func (t *totaler) percentSolutions(value int64) float64 {
	return fw.Percent(value, t.result.solutions)
}
