package cmd

import (
	"fmt"
	"github.com/akutz/sortfold"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"io"
	"runtime"
	"strings"
)

func normalize(s string) string {
	return strings.ToUpper(s)
}

func outputSortedMap(writer io.Writer, itemsMap map[string][]string, keyPrefix string) {
	var keys []string
	for k := range itemsMap {
		keys = append(keys, k)
	}

	sortfold.Strings(keys)

	for _, k := range keys {
		color.Fprintf(writer, "\n<gray>%s: %s</>\n", keyPrefix, k)
		sortAndOutput(writer, itemsMap[k])
	}
}

func sortAndOutput(writer io.Writer, items []string) {
	sortfold.Strings(items)
	for _, item := range items {
		_, _ = fmt.Fprintf(writer, " %s\n", item)
	}
}

// printMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func printMemUsage(w io.Writer) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	color.Fprintf(w, "\n<gray>Alloc =</> <green>%s</>", humanize.IBytes(m.Alloc))
	color.Fprintf(w, "\t<gray>TotalAlloc =</> <green>%s</>", humanize.IBytes(m.TotalAlloc))
	color.Fprintf(w, "\t<gray>Sys =</> <green>%s</>", humanize.IBytes(m.Sys))
	color.Fprintf(w, "\t<gray>NumGC =</> <green>%v</>\n", m.NumGC)
}
