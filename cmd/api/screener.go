package api

import (
	"github.com/akutz/sortfold"
)

type Screener struct {
	p Printer
}

func NewScreener(p Printer) *Screener {
	s := Screener{
		p: p,
	}
	return &s
}

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

func (s *Screener) WriteSlice(items []string) {
	sortfold.Strings(items)
	for _, item := range items {
		s.p.Cprint(" %s\n", item)
	}
}
