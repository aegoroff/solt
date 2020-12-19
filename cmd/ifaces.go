package cmd

import (
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/gookit/color"
	"github.com/spf13/afero"
	"io"
	"text/tabwriter"
)

// Matcher defines string matcher interface
type Matcher interface {
	// Match do string matching to several patterns
	Match(s string) bool
}

type nugetprinter interface {
	print(parent string, packs []*pack)
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

type screener interface {
	writeMap(m map[string][]string, keyPrefix string)
	writeSlice(slice []string)
}

type conf interface {
	// fs defines app file system abstraction
	fs() afero.Fs

	prn() printer

	sourcesPath() *string
}

type command interface {
	execute() error
}

type sdkModuleHandler interface {
	handle(sol string, refs map[string]c9s.StringHashSet)
}
