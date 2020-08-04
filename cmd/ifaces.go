package cmd

// Matcher defines string matcher interface
type Matcher interface {
	// Match do string matching to several patterns
	Match(s string) bool
}

type nugetprinter interface {
	print(parent string, packs []*pack)
}
