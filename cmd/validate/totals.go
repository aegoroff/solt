package validate

import (
	"solt/internal/out"
	"solt/internal/ux"
)

type totals struct {
	solutions        int64
	problemSolutions int64
	projects         int64
	problemProjects  int64
	redundantRefs    int64
}

func (t *totals) display(p out.Printer, w out.Writable) {
	p.Println()
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)

	tbl.AddHead("Parameter", "Count", "%     ")

	sl := ux.NewLine("Solutions", t.solutions)
	tbl.AddLine(sl.Name(), sl.Value(), "")

	psl := ux.NewLine("Problem solutions", t.problemSolutions)
	tbl.AddLine(psl.Name(), psl.Value(), psl.Percent(t.solutions))

	pl := ux.NewLine("SDK Projects", t.projects)
	tbl.AddLine(pl.Name(), pl.Value(), "")

	ppl := ux.NewLine("Problem projects", t.problemProjects)
	tbl.AddLine(ppl.Name(), ppl.Value(), ppl.Percent(t.projects))

	tbl.AddLines(ux.NewLine("Redundant references", t.redundantRefs))
	tbl.Print()
}
