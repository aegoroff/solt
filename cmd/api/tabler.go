package api

import "github.com/cheynewallace/tabby"

// Tabler table drawing component
type Tabler struct {
	margin *Marginer
	prn    Printer
	tab    *tabby.Tabby
}

// NewTabler creates new Tabler instance
func NewTabler(prn Printer, margin int) *Tabler {
	return &Tabler{
		prn:    prn,
		margin: NewMarginer(margin),
		tab:    tabby.NewCustom(prn.Twriter()),
	}
}

// AddHead adds headline
func (t *Tabler) AddHead(columns ...string) {
	names := make([]interface{}, len(columns))
	underlines := make([]interface{}, len(columns))
	for i, column := range columns {
		under := newUnderline(column)
		if i == 0 {
			underlines[i] = t.margin.Margin(under)
			names[i] = t.margin.Margin(column)
		} else {
			underlines[i] = under
			names[i] = column
		}
	}
	t.tab.AddLine(names...)
	t.tab.AddLine(underlines...)
}

// AddLine adds data line to the table
func (t Tabler) AddLine(line ...string) {
	data := make([]interface{}, len(line))
	for i, column := range line {
		if i == 0 {
			data[i] = t.margin.Margin(column)
		} else {
			data[i] = column
		}
	}

	t.tab.AddLine(data...)
}

// Print prints table
func (t *Tabler) Print() {
	t.tab.Print()
}

// newUnderline creates new dashes string that can be used as underline
// the length of the line is equal len of string specified
func newUnderline(s string) string {
	return NewCustomMarginer(len(s), "-").Margin("")
}
