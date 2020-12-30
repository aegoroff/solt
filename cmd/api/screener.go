package api

import (
	"github.com/akutz/sortfold"
)

// Screener is an abstraction that does comples structures output
type Screener struct {
	p Printer
	m *Marginer
}

// NewScreener creates new Screener instance
func NewScreener(p Printer) *Screener {
	s := Screener{
		p: p,
		m: NewMarginer(1),
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
		s.p.Cprint("\n<gray>%s: %s</>\n", keyPrefix, k)
		s.WriteSlice(itemsMap[k])
	}
}

// WriteSlice prints []string instance
func (s *Screener) WriteSlice(items []string) {
	sortfold.Strings(items)
	for _, item := range items {
		s.p.Cprint(s.m.Margin("%s\n"), item)
	}
}
