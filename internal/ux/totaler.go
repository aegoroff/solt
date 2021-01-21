package ux

import "solt/internal/out"

// Totaler does total in output
type Totaler struct {
	p out.Printer
	w out.Writable
}

// NewTotaler creates new Totaler instance
func NewTotaler(p out.Printer, w out.Writable) *Totaler {
	return &Totaler{p: p, w: w}
}

// Display outputs total in
func (t *Totaler) Display(d Displayer) {
	t.p.Cprint(" <red>Totals:</>")
	t.p.Println()
	t.p.Println()
	tbl := NewTabler(t.w, 2)
	d.Display(tbl)
	tbl.Print()
}
