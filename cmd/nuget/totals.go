package nuget

import (
	"github.com/dustin/go-humanize"
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
	tbl.AddLine("Solutions", humanize.Comma(t.solutions))
	tbl.AddLine("", "")
	tbl.AddHead("Packages", "Count")

	nl := ux.NewLine("Total", t.nugets)

	tbl.AddLine(nl.Name(), nl.Value())

	ml := ux.NewLine("Mismatched", t.mismatched)
	tbl.AddLine(ml.Name(), ml.Value())

	tbl.Print()
}

func (t *totalsByProjects) display(p out.Printer, w out.Writable) {
	p.Println()
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)
	tbl.AddLine("Projects", humanize.Comma(t.projects))
	tbl.AddLine("Packages", humanize.Comma(t.nugets))

	tbl.Print()
}
