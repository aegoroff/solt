package cmd

import (
	"github.com/gookit/color"
	"io"
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
