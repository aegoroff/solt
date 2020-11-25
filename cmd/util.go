package cmd

import (
	"github.com/akutz/sortfold"
	"github.com/dustin/go-humanize"
	"runtime"
	"strings"
)

func normalize(s string) string {
	return strings.ToUpper(s)
}

type screenerImpl struct {
	p printer
}

func newScreener(p printer) screener {
	s := screenerImpl{
		p: p,
	}
	return &s
}

func (s *screenerImpl) writeMap(itemsMap map[string][]string, keyPrefix string) {
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

func (s *screenerImpl) writeSlice(items []string) {
	sortfold.Strings(items)
	for _, item := range items {
		s.p.cprint(" %s\n", item)
	}
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(p printer) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	p.cprint("\n<gray>Alloc =</> <green>%s</>", humanize.IBytes(m.Alloc))
	p.cprint("\t<gray>TotalAlloc =</> <green>%s</>", humanize.IBytes(m.TotalAlloc))
	p.cprint("\t<gray>Sys =</> <green>%s</>", humanize.IBytes(m.Sys))
	p.cprint("\t<gray>NumGC =</> <green>%v</>\n", m.NumGC)
}
