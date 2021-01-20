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
	tbl.AddLine(pl.Name(), pl.Value())
	tbl.AddLine("", "")

	tbl.AddHead("Files", "Count")

	lines := ux.NewLines()
	lines.Add("Found", t.found)
	lines.Add("Included", t.included)
	lines.Add("Lost", t.lost)
	lines.Add("Included but not exist", t.unexist)

	for _, l := range lines {
		tbl.AddLine(l.Name(), l.Value())
	}

	tbl.Print()
}
