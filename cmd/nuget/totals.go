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
	tbl.AddStringLine("Solutions", humanize.Comma(t.solutions))
	tbl.AddStringLine("", "")
	tbl.AddHead("Packages", "Count")

	nl := ux.NewLine("Total", t.nugets)
	ml := ux.NewLine("Mismatched", t.mismatched)

	tbl.AddLines(nl, ml)

	tbl.Print()
}

func (t *totalsByProjects) display(p out.Printer, w out.Writable) {
	p.Println()
	p.Cprint(" <red>Totals:</>\n\n")

	tbl := ux.NewTabler(w, 2)
	tbl.AddStringLine("Projects", humanize.Comma(t.projects))
	tbl.AddStringLine("Packages", humanize.Comma(t.nugets))

	tbl.Print()
}
