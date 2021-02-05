package lf

import (
	"solt/internal/fw"
	"strings"
)

type finder struct {
	m fw.Matcher
}

func newFinder(includedFiles, excludedFolders []string) *finder {
	excludes := exclude(excludedFolders)
	includes := fw.NewExactMatch(includedFiles, false)

	m := fw.NewLostItemMatcher(includes, excludes)
	return &finder{m: m}
}

func exclude(excludedFolders []string) fw.Matcher {
	excludes, err := fw.NewPartialMatcher(excludedFolders, strings.ToUpper)
	if err != nil {
		return fw.NewMatchNothing()
	}
	return excludes
}

func (l *finder) find(files []string) []string {
	return fw.Filter(files, l.m)
}
