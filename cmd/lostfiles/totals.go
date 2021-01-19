package lostfiles

import (
	"github.com/dustin/go-humanize"
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
	tbl.AddLine("Projects", humanize.Comma(t.projects))
	tbl.AddLine("", "")

	tbl.AddHead("Files", "Count")

	type hv struct {
		h string
		v int64
	}

	lines := []hv{
		{"Found", t.found},
		{"Included", t.included},
		{"Lost", t.lost},
		{"Included but not exist", t.unexist},
	}
	for _, l := range lines {
		line := newTotalLine(l.h, l.v)
		tbl.AddLine(line.h(), line.v())
	}

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
