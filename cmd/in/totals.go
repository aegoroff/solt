package in

import (
	"fmt"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/dustin/go-humanize"
	"solt/internal/ux"
	"strings"
)

type totals struct {
	solutions    int64
	projects     int64
	projectTypes rbtree.RbTree
}

type typeStat struct {
	name      string
	count     int64
	solutions int64
}

func (t *totals) Display(tbl *ux.Tabler) {
	tbl.AddLines(
		ux.NewLine("Solutions", t.solutions),
		ux.NewLine("Projects", t.projects),
	)
	tbl.AddLine("", "")

	const percentH = "%     "
	tbl.AddHead("Project type", "Count", percentH, "Solutions", percentH)

	it := rbtree.NewAscend(t.projectTypes)

	it.Foreach(func(n rbtree.Comparable) {
		ts := n.(*typeStat)
		tbl.AddLine(
			ts.name,
			ts.Count(),
			ts.CountPercent(t.projects),
			ts.Solutions(),
			ts.SolutionsPercent(t.solutions),
		)
	})
}

func newTypeStat(name string) *typeStat {
	return &typeStat{name: name}
}

func (t *typeStat) Less(y rbtree.Comparable) bool {
	return strings.Compare(t.name, y.(*typeStat).name) < 0
}

func (t *typeStat) Equal(y rbtree.Comparable) bool {
	return strings.Compare(t.name, y.(*typeStat).name) == 0
}

func (t *typeStat) Count() string {
	return humanize.Comma(t.count)
}

func (t *typeStat) CountPercent(tot int64) string {
	p := ux.Percent(t.count, tot)
	return fmt.Sprintf("%.2f%%", p)
}

func (t *typeStat) Solutions() string {
	return humanize.Comma(t.solutions)
}

func (t *typeStat) SolutionsPercent(tot int64) string {
	p := ux.Percent(t.solutions, tot)
	return fmt.Sprintf("%.2f%%", p)
}
