package cmd

import (
	"bytes"
	c9s "github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"github.com/anknown/ahocorasick"
	"solt/msvc"
)

// matchP defines partial matching
type matchP struct {
	machine *goahocorasick.Machine
}

// matchE defines exact matching
type matchE struct {
	hashset c9s.StringHashSet
}

// matchTree defines exact matching using rbtree.RbTree
type matchTree struct {
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
	include   Matcher
	exclude   Matcher
	decorator msvc.StringDecorator
}

// NewLostItemMatcher creates new Matcher instance that detects lost item
func NewLostItemMatcher(incl Matcher, excl Matcher, decorator msvc.StringDecorator) Matcher {
	return &matchL{
		include:   incl,
		exclude:   excl,
		decorator: decorator,
	}
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

// NewExactMatch creates exacth matcher from strings slice
func NewExactMatch(matches []string, decorator msvc.StringDecorator) Matcher {
	h := make(c9s.StringHashSet, len(matches))
	for _, s := range matches {
		h.Add(decorator(s))
	}

	hs := matchE{hashset: h}
	return &hs
}

// NewExactTreeMatch creates exacth matcher from strings slice
func NewExactTreeMatch(matches []string) Matcher {
	tree := rbtree.NewRbTree()
	for _, s := range matches {
		cs := caseless(s)
		tree.Insert(&cs)
	}

	hs := matchTree{tree: tree}
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

func (m *matchL) Match(s string) bool {
	decorated := m.decorator(s)
	return !m.include.Match(decorated) && !m.exclude.Match(decorated)
}

func (m *matchTree) Match(s string) bool {
	cs := caseless(s)
	_, ok := m.tree.Search(&cs)
	return ok
}
