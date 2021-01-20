package nuget

import (
	"solt/internal/out"
	"solt/internal/ux"
)

type totalsBySolution struct {
	solutions  int64
	nugets     int64
	mismatched int64
}

type totalsByProjects struct {
	projects int64
	nugets   int64
}

func (t *totalsBySolution) display(p out.Printer, w out.Writable) {
	p.Println()
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)
	tbl.AddLines(ux.NewLine("Solutions", t.solutions))
	tbl.AddLine("", "")
	tbl.AddHead("Packages", "Count")

	tbl.AddLines(
		ux.NewLine("Total", t.nugets),
		ux.NewLine("Mismatched", t.mismatched),
	)

	tbl.Print()
}

func (t *totalsByProjects) display(p out.Printer, w out.Writable) {
	p.Println()
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)

	tbl.AddLines(
		ux.NewLine("Projects", t.projects),
		ux.NewLine("Packages", t.nugets),
	)

	tbl.Print()
}
