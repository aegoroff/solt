package lostprojects

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"solt/cmd/fw"
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
	tbl.AddLine("Solutions", humanize.Comma(t.solutions))
	tbl.AddLine("Projects", humanize.Comma(t.allProjects))
	tbl.AddLine("", "")

	within := t.allProjects - t.lost - t.lostWithIncludes
	tbl.AddHead("", "Count", "%     ")

	type kv struct {
		h string
		v int64
	}

	lines := []kv{
		{"Within solutions", within},
		{"Lost projects", t.lost},
		{"Lost projects with includes", t.lostWithIncludes},
		{"Included but not exist", t.unexist},
		{"Removed (if specified)", t.removed},
	}
	for _, l := range lines {
		line := newTotalLine(l.h, l.v, t.allProjects)
		tbl.AddLine(line.h(), line.v(), line.p())
	}

	tbl.Print()
}

type totalLine struct {
	head    string
	val     int64
	percent float64
}

func (t *totalLine) h() string { return t.head }
func (t *totalLine) v() string { return humanize.Comma(t.val) }
func (t *totalLine) p() string { return fmt.Sprintf("%.2f%%", t.percent) }

func newTotalLine(head string, val int64, tot int64) *totalLine {
	return &totalLine{
		head:    head,
		val:     val,
		percent: fw.Percent(val, tot),
	}
}
