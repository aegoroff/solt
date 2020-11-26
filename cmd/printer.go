package cmd

import (
	"fmt"
	"github.com/akutz/sortfold"
	"github.com/gookit/color"
	"io"
	"text/tabwriter"
)

type prn struct {
	tw *tabwriter.Writer
	w  io.Writer
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

// newPrinter creates new printer interface instance
func newPrinter(w io.Writer) printer {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 4, ' ', 0)

	p := prn{
		tw: tw,
		w:  w,
	}
	return &p
}

func (r *prn) writer() io.Writer {
	return r.w
}

func (r *prn) twriter() *tabwriter.Writer {
	return r.tw
}

func (r *prn) flush() {
	_ = r.tw.Flush()
}

func (r *prn) tprint(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(r.tw, format, a...)
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

func (r *prn) cprint(format string, a ...interface{}) {
	color.Fprintf(r.w, format, a...)
}

func (*prn) setColor(c color.Color) {
	_, _ = color.Set(c)
}

func (*prn) resetColor() {
	_, _ = color.Reset()
}
