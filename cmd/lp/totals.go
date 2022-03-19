package lp

import (
	"solt/internal/ux"
)

type totals struct {
	solutions        int64
	allProjects      int64
	lost             int64
	lostWithIncludes int64
	missing          int64
	removed          int64
}

func (t *totals) Display(tbl *ux.Tabler) {
	tbl.AddLines(
		ux.NewLine("Solutions", t.solutions),
		ux.NewLine("Projects", t.allProjects),
	)
	tbl.AddLine("", "")

	within := t.allProjects - t.lost - t.lostWithIncludes
	tbl.AddHead("", "Count", "%     ")

	lines := ux.NewLines()
	lines.Add("Within solutions", within)
	lines.Add("Lost projects", t.lost)
	lines.Add("Lost projects with includes", t.lostWithIncludes)
	lines.Add("Included but not exist", t.missing)
	lines.Add("Removed (if specified)", t.removed)

	for _, l := range lines {
		tbl.AddLine(l.Name(), l.Value(), l.Percent(t.allProjects))
	}
}
