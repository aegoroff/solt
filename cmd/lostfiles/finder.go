package lostfiles

import (
	"solt/cmd/fw"
	"strings"
)

type finder struct {
	m fw.Matcher
}

func newFinder(includedFiles, excludedFolders []string) (*finder, error) {
	excludes, err := fw.NewPartialMatcher(excludedFolders, strings.ToUpper)
	if err != nil {
		return nil, err
	}

	includes := fw.NewExactMatch(includedFiles)

	m := fw.NewLostItemMatcher(includes, excludes)

	return &finder{m: m}, nil
}

func (l *finder) find(files []string) []string {
	return fw.Filter(files, l.m)
}
