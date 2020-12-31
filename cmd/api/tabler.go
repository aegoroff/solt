package api

import (
	"github.com/cheynewallace/tabby"
)

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
func (t *Tabler) AddHead(line ...string) {
	t.addLine(line, originalString, underline)
}

// AddLine adds data line to the table
func (t Tabler) AddLine(line ...string) {
	t.addLine(line, originalString)
}

func (t Tabler) addLine(line []string, decors ...func(s string) string) {
	for _, d := range decors {
		t.tab.AddLine(t.newLine(line, d)...)
	}
}

func (t Tabler) newLine(line []string, d func(s string) string) []interface{} {
	result := make([]interface{}, len(line))
	for i, column := range line {
		decorated := d(column)
		if i == 0 {
			decorated = t.margin.Margin(decorated)
		}
		result[i] = decorated
	}

	return result
}

// Print prints table
func (t *Tabler) Print() {
	t.tab.Print()
}

func originalString(s string) string { return s }

// underline creates new dashes string that can be used as underline
// the length of the line is equal len of string specified
func underline(s string) string {
	return NewCustomMarginer(len(s), "-").Margin("")
}
