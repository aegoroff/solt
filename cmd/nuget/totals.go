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

	nl := newTotalLine("Total", t.nugets)
	tbl.AddLine(nl.h(), nl.v())

	ml := newTotalLine("Mismatched", t.mismatched)
	tbl.AddLine(ml.h(), ml.v())

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

type totalLine struct {
	head string
	val  int64
}

func (t *totalLine) h() string { return t.head }
func (t *totalLine) v() string { return humanize.Comma(t.val) }

func newTotalLine(head string, val int64) *totalLine {
	return &totalLine{
		head: head,
		val:  val,
	}
}
