package ux

import (
	"github.com/cheynewallace/tabby"
	"solt/internal/out"
	"text/tabwriter"
)

// Tabler table drawing component
type Tabler struct {
	margin *Marginer
	tab    *tabby.Tabby
}

// NewTabler creates new Tabler instance
func NewTabler(w out.Writable, margin int) *Tabler {
	tw := new(tabwriter.Writer).Init(w.Writer(), 0, 8, 4, ' ', 0)
	return &Tabler{
		margin: NewMarginer(margin),
		tab:    tabby.NewCustom(tw),
	}
}

// AddHead adds headline
func (t *Tabler) AddHead(line ...string) {
	t.addLine(line, originalString, underline)
}

// AddStringLine adds data line to the table
func (t Tabler) AddStringLine(line ...string) {
	t.addLine(line, originalString)
}

// AddLine adds new *Line into table
func (t Tabler) AddLine(line *Line) {
	t.AddStringLine(line.Name(), line.Value())
}

// AddLines adds many lines
func (t Tabler) AddLines(lines ...*Line) {
	for _, line := range lines {
		t.AddLine(line)
	}
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
