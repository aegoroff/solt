package fw

import (
	"bytes"
	"github.com/akutz/sortfold"
	"github.com/anknown/ahocorasick"
	"github.com/google/btree"
)

// matchP defines partial matching
type matchP struct {
	machine   *goahocorasick.Machine
	decorator func(s string) string
}

// matchE defines exact matching using btree.BTree
type matchE struct {
	tree *btree.BTree
}

// matchL defines lost item matching
type matchL struct {
	include Matcher
	exclude Matcher
}

type caseless string

func (c *caseless) Less(than btree.Item) bool {
	return sortfold.CompareFold(string(*c), string(*than.(*caseless))) < 0
}

type matchNothing struct{}

// NewMatchNothing creates new Matcher that matches nothing i.e. always return false on match
func NewMatchNothing() Matcher {
	return &matchNothing{}
}

func (*matchNothing) Match(string) bool { return false }

type matchAll struct {
	matchers []Matcher
}

// NewMatchAll creates new Matcher that matches when all matchers match
func NewMatchAll(matchers ...Matcher) Matcher {
	return &matchAll{
		matchers: matchers,
	}
}

func (ma *matchAll) Match(s string) bool {
	for _, matcher := range ma.matchers {
		if !matcher.Match(s) {
			return false
		}
	}
	return true
}

// NewLostItemMatcher creates new Matcher instance that detects lost item
func NewLostItemMatcher(incl Matcher, excl Matcher) Matcher {
	m := &matchL{
		include: incl,
		exclude: excl,
	}
	return m
}

// NewPartialMatcher creates new matcher that implements Aho corasick multi pattern matching
// Partial means that string should contain one of the matcher's strings as substring
// or whole string
func NewPartialMatcher(matches []string, decorator func(s string) string) (Searcher, error) {
	runes := make([][]rune, len(matches))
	for i, s := range matches {
		ds := decorator(s)
		runes[i] = bytes.Runes([]byte(ds))
	}
	machine := new(goahocorasick.Machine)
	err := machine.Build(runes)
	if err != nil {
		return nil, err
	}
	aho := &matchP{
		machine:   machine,
		decorator: decorator,
	}
	return aho, nil
}

// NewExactMatch creates exact matcher from strings slice
func NewExactMatch(matches []string) Matcher {
	tree := btree.New(1024)

	for _, s := range matches {
		cs := caseless(s)
		tree.ReplaceOrInsert(&cs)
	}

	m := &matchE{
		tree: tree,
	}
	return m
}

// Match do string matching to several patterns
func (a *matchP) Match(s string) bool {
	ds := a.decorator(s)
	terms := a.machine.MultiPatternSearch([]rune(ds), true)
	return len(terms) > 0
}

func (a *matchP) Search(s string) []string {
	ds := a.decorator(s)
	terms := a.machine.MultiPatternSearch([]rune(ds), false)
	result := make([]string, len(terms))
	for i, term := range terms {
		result[i] = string(term.Word)
	}
	return result
}

func (m *matchL) Match(s string) bool {
	return !m.include.Match(s) && !m.exclude.Match(s)
}

func (m *matchE) Match(s string) bool {
	cs := caseless(s)
	return m.tree.Get(&cs) != nil
}

// MatchAny does any string matching to several patterns
func MatchAny(ss []string, m Matcher) bool {
	for _, s := range ss {
		if m.Match(s) {
			return true
		}
	}

	return false
}

// Filter filters slice using Matcher. Only matched strings will be in result
// IMPORTANT: source slice MUST NOT be used after calling this method
func Filter(ss []string, m Matcher) []string {
	if m == nil {
		return []string{}
	}
	n := 0
	for _, s := range ss {
		if m.Match(s) {
			ss[n] = s
			n++
		}
	}
	ss = ss[:n]
	return ss
}
