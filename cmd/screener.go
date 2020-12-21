package cmd

import "github.com/akutz/sortfold"

type screener struct {
	p printer
}

func newScreener(p printer) *screener {
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
		s.p.cprint("\n<gray>%s: %s</>\n", keyPrefix, k)
		s.writeSlice(itemsMap[k])
	}
}

func (s *screener) writeSlice(items []string) {
	sortfold.Strings(items)
	for _, item := range items {
		s.p.cprint(" %s\n", item)
	}
}
