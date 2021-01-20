package lostfiles

import (
	"solt/internal/out"
	"solt/internal/ux"
)

type totals struct {
	projects int64
	unexist  int64
	included int64
	lost     int64
	found    int64
}

func (t *totals) display(p out.Printer, w out.Writable) {
	p.Println()
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)
	pl := ux.NewLine("Projects", t.projects)
	tbl.AddLine(pl)
	tbl.AddStringLine("", "")

	tbl.AddHead("Files", "Count")

	lines := ux.NewLines()
	lines.Add("Found", t.found)
	lines.Add("Included", t.included)
	lines.Add("Lost", t.lost)
	lines.Add("Included but not exist", t.unexist)

	tbl.AddLines(lines...)

	tbl.Print()
}
