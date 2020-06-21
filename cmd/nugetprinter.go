package cmd

import (
	"fmt"
	"github.com/akutz/sortfold"
	"github.com/gookit/color"
	"io"
	"sort"
	"strings"
	"text/tabwriter"
)

// pack defines nuget package descriptor
type pack struct {
	pkg      string
	versions []string
}

func newNugetPrinter(w io.Writer) nugetprinter {
	p := nugetprint{
		w:  w,
		tw: new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0),
	}
	return &p
}

type nugetprinter interface {
	print(parent string, packs []*pack)
}

type nugetprint struct {
	w  io.Writer
	tw *tabwriter.Writer
}

func (n *nugetprint) print(parent string, packs []*pack) {
	color.Fprintf(n.w, "\n <gray>%s</>\n", parent)

	const format = "  %v\t%v\n"
	_, _ = fmt.Fprintf(n.tw, format, "Package", "Version")
	_, _ = fmt.Fprintf(n.tw, format, "-------", "-------")

	sort.Slice(packs, func(i, j int) bool {
		return sortfold.CompareFold(packs[i].pkg, packs[j].pkg) < 0
	})

	for _, item := range packs {
		_, _ = fmt.Fprintf(n.tw, format, item.pkg, strings.Join(item.versions, ", "))
	}
	_ = n.tw.Flush()
}
