package ux

import (
	"fmt"
	"github.com/dustin/go-humanize"
)

// Line provides statistic line
type Line struct {
	name string
	val  int64
}

// Lines provides *Line slice
type Lines []*Line

// NewLines creates new *Line slice
func NewLines() Lines {
	return make(Lines, 0)
}

// Add adds new line to lines slice
func (s *Lines) Add(name string, val int64) {
	*s = append(*s, NewLine(name, val))
}

// Name gets a parameter name
func (t *Line) Name() string { return t.name }

// Value gets parameter value
func (t *Line) Value() string { return humanize.Comma(t.val) }

// Percent calculates value percent
func (t *Line) Percent(total int64) string {
	pv := Percent(t.val, total)
	return fmt.Sprintf("%.2f%%", pv)
}

// NewLine creates new *Line
func NewLine(name string, val int64) *Line {
	return &Line{
		name: name,
		val:  val,
	}
}
