package lostfiles

import (
	"solt/cmd/fw"
	"strings"
)

type finder struct {
	m fw.Matcher
}

type matchNothing struct{}

func (*matchNothing) Match(string) bool { return false }

func newFinder(includedFiles, excludedFolders []string) *finder {
	excludes := exclude(excludedFolders)
	includes := fw.NewExactMatch(includedFiles)

	m := fw.NewLostItemMatcher(includes, excludes)
	return &finder{m: m}
}

func exclude(excludedFolders []string) fw.Matcher {
	excludes, err := fw.NewPartialMatcher(excludedFolders, strings.ToUpper)
	if err != nil {
		return &matchNothing{}
	}
	return excludes
}

func (l *finder) find(files []string) []string {
	return fw.Filter(files, l.m)
}
