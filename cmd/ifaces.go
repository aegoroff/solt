package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/gookit/color"
	"io"
	"text/tabwriter"
)

// Matcher defines string matcher interface
type Matcher interface {
	// Match do string matching to several patterns
	Match(s string) bool
}

type printer interface {
	writer() io.Writer

	twriter() *tabwriter.Writer

	flush()

	// tprint prints using tab writer
	tprint(format string, a ...interface{})

	// cprint prints data with suppport colorizing
	cprint(format string, a ...interface{})

	// setColor sets console color
	setColor(c color.Color)

	// resetColor resets console color
	resetColor()
}

type executor interface {
	execute() error
}

type sdkActioner interface {
	action(solution string, items map[string]c9s.StringHashSet)
}
