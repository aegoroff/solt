package lf

import (
	"solt/internal/ux"
)

type totals struct {
	projects int64
	missing  int64
	included int64
	lost     int64
	found    int64
}

func (t *totals) Display(tbl *ux.Tabler) {
	pl := ux.NewLine("Projects", t.projects)
	tbl.AddLines(pl)
	tbl.AddLine("", "")

	tbl.AddHead("Files", "Count")

	lines := ux.NewLines()
	lines.Add("Found", t.found)
	lines.Add("Included", t.included)
	lines.Add("Included but not exist", t.missing)
	lines.Add("Lost", t.lost)

	tbl.AddLines(lines...)
}
