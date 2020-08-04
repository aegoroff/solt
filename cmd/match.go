package cmd

import (
	"bytes"
	"github.com/aegoroff/godatastruct/collections"
	"github.com/anknown/ahocorasick"
)

// matchP defines partial matching
type matchP struct {
	machine *goahocorasick.Machine
}

// matchP defines exact matching
type matchE struct {
	hashset collections.StringHashSet
}

// NewPartialMatcher creates new matcher that implements Aho corasick multi pattern matching
// Partial means that string should contain one of the matcher's strings as substring
// or whole string
func NewPartialMatcher(matches []string) (Matcher, error) {
	runes := make([][]rune, 0, len(matches))
	for _, s := range matches {
		runes = append(runes, bytes.Runes([]byte(s)))
	}
	machine := new(goahocorasick.Machine)
	err := machine.Build(runes)
	if err != nil {
		return nil, err
	}
	aho := matchP{machine: machine}

	return &aho, nil
}

// NewExactMatchS creates exacth matcher from strings slice
// Exact means that string must exactly match one of the matcher's strings
func NewExactMatchS(matches []string) Matcher {
	h := make(collections.StringHashSet, len(matches))
	for _, s := range matches {
		h.Add(s)
	}

	return NewExactMatchHS(&h)
}

// NewExactMatchHS creates exacth matcher from strings hashset
// Exact means that string must exactly match one of the matcher's strings
func NewExactMatchHS(existing *collections.StringHashSet) Matcher {
	hs := matchE{hashset: *existing}
	return &hs
}

// Match do string matching to several patterns
func (a *matchP) Match(s string) bool {
	terms := a.machine.MultiPatternSearch([]rune(s), true)
	return len(terms) > 0
}

// Match do string matching to several patterns
func (h *matchE) Match(s string) bool {
	return h.hashset.Contains(s)
}
