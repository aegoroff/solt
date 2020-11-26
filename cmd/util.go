package cmd

import (
	"github.com/dustin/go-humanize"
	"runtime"
	"strings"
)

func normalize(s string) string {
	return strings.ToUpper(s)
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
