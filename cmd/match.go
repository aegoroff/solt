package cmd

import (
	"bytes"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"github.com/anknown/ahocorasick"
	"solt/msvc"
)

// matchP defines partial matching
type matchP struct {
	machine   *goahocorasick.Machine
	decorator msvc.StringDecorator
}

// matchE defines exact matching using rbtree.RbTree
type matchE struct {
	tree rbtree.RbTree
}

type caseless string

func (c *caseless) LessThan(y rbtree.Comparable) bool {
	return c.compare(y) < 0
}

func (c *caseless) EqualTo(y rbtree.Comparable) bool {
	return c.compare(y) == 0
}

func (c *caseless) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(string(*c), string(*y.(*caseless)))
}

type matchL struct {
	include Matcher
	exclude Matcher
}

// NewLostItemMatcher creates new Matcher instance that detects lost item
func NewLostItemMatcher(incl Matcher, excl Matcher) Matcher {
	return &matchL{
		include: incl,
		exclude: excl,
	}
}

// NewPartialMatcher creates new matcher that implements Aho corasick multi pattern matching
// Partial means that string should contain one of the matcher's strings as substring
// or whole string
func NewPartialMatcher(matches []string, decorator msvc.StringDecorator) (Matcher, error) {
	runes := make([][]rune, 0, len(matches))
	for _, s := range matches {
		ds := decorator(s)
		runes = append(runes, bytes.Runes([]byte(ds)))
	}
	machine := new(goahocorasick.Machine)
	err := machine.Build(runes)
	if err != nil {
		return nil, err
	}
	aho := matchP{
		machine:   machine,
		decorator: decorator,
	}

	return &aho, nil
}

// NewExactMatch creates exact matcher from strings slice
func NewExactMatch(matches []string) Matcher {
	tree := rbtree.NewRbTree()
	for _, s := range matches {
		cs := caseless(s)
		tree.Insert(&cs)
	}

	hs := matchE{tree: tree}
	return &hs
}

// Match do string matching to several patterns
func (a *matchP) Match(s string) bool {
	ds := a.decorator(s)
	terms := a.machine.MultiPatternSearch([]rune(ds), true)
	return len(terms) > 0
}

func (m *matchL) Match(s string) bool {
	return !m.include.Match(s) && !m.exclude.Match(s)
}

func (m *matchE) Match(s string) bool {
	cs := caseless(s)
	_, ok := m.tree.SearchNode(&cs)
	return ok
}
