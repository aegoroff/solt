package ux

import (
	"github.com/akutz/sortfold"
	"solt/internal/out"
)

// Screener is an abstraction that does comples structures output
type Screener struct {
	p  out.Printer
	m1 *Marginer
	m2 *Marginer
}

// NewScreener creates new Screener instance
func NewScreener(p out.Printer) *Screener {
	s := Screener{
		p:  p,
		m1: NewMarginer(1),
		m2: NewMarginer(2),
	}
	return &s
}

// WriteMap prints map[string][]string instance
func (s *Screener) WriteMap(itemsMap map[string][]string, keyPrefix string) {
	var keys []string
	for k := range itemsMap {
		keys = append(keys, k)
	}

	sortfold.Strings(keys)

	for _, k := range keys {
		s.p.Println()
		s.p.Cprint(s.m1.Margin("<gray>%s: %s</>\n"), keyPrefix, k)
		s.WriteSlice(itemsMap[k])
	}
}

// WriteSlice prints []string instance
func (s *Screener) WriteSlice(items []string) {
	sortfold.Strings(items)
	for _, item := range items {
		s.p.Cprint(s.m2.Margin("%s\n"), item)
	}
}
