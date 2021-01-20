package lostprojects

import (
	"solt/internal/out"
	"solt/internal/ux"
)

type totals struct {
	solutions        int64
	allProjects      int64
	lost             int64
	lostWithIncludes int64
	unexist          int64
	removed          int64
}

func (t *totals) display(p out.Printer, w out.Writable) {
	p.Println()
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)
	sl := ux.NewLine("Solutions", t.solutions)
	pl := ux.NewLine("Projects", t.allProjects)
	tbl.AddLines(sl, pl)
	tbl.AddStringLine("", "")

	within := t.allProjects - t.lost - t.lostWithIncludes
	tbl.AddHead("", "Count", "%     ")

	lines := ux.NewLines()
	lines.Add("Within solutions", within)
	lines.Add("Lost projects", t.lost)
	lines.Add("Lost projects with includes", t.lostWithIncludes)
	lines.Add("Included but not exist", t.unexist)
	lines.Add("Removed (if specified)", t.removed)

	for _, l := range lines {
		tbl.AddStringLine(l.Name(), l.Value(), l.Percent(t.allProjects))
	}

	tbl.Print()
}
