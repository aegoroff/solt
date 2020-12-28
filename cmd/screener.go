package cmd

import (
	"github.com/akutz/sortfold"
	"solt/cmd/api"
)

type screener struct {
	p api.Printer
}

func newScreener(p api.Printer) *screener {
	s := screener{
		p: p,
	}
	return &s
}

func (s *screener) writeMap(itemsMap map[string][]string, keyPrefix string) {
	var keys []string
	for k := range itemsMap {
		keys = append(keys, k)
	}

	sortfold.Strings(keys)

	for _, k := range keys {
		s.p.Cprint("\n<gray>%s: %s</>\n", keyPrefix, k)
		s.writeSlice(itemsMap[k])
	}
}

func (s *screener) writeSlice(items []string) {
	sortfold.Strings(items)
	for _, item := range items {
		s.p.Cprint(" %s\n", item)
	}
}
