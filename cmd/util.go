package cmd

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/gookit/color"
	"io"
	"runtime"
	"sort"
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

	sort.Strings(keys)

	for _, k := range keys {
		color.Fprintf(writer, "\n<gray>%s: %s</>\n", keyPrefix, k)
		sortAndOutput(writer, itemsMap[k])
	}
}

func sortAndOutput(writer io.Writer, items []string) {
	sort.Strings(items)
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
	_, _ = fmt.Fprintf(w, "\nAlloc = %s", humanize.IBytes(m.Alloc))
	_, _ = fmt.Fprintf(w, "\tTotalAlloc = %s", humanize.IBytes(m.TotalAlloc))
	_, _ = fmt.Fprintf(w, "\tSys = %s", humanize.IBytes(m.Sys))
	_, _ = fmt.Fprintf(w, "\tNumGC = %v\n", m.NumGC)
}
