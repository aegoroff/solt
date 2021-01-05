package lostfiles

import (
	"solt/cmd/api"
	"strings"
)

type lostFinder struct {
	m api.Matcher
}

func newLostFinder(includedFiles, excludedFolders []string) (*lostFinder, error) {
	excludes, err := api.NewPartialMatcher(excludedFolders, strings.ToUpper)
	if err != nil {
		return nil, err
	}

	includes := api.NewExactMatch(includedFiles)

	m := api.NewLostItemMatcher(includes, excludes)

	return &lostFinder{m: m}, nil
}

func (l *lostFinder) find(files []string) []string {
	return api.Filter(files, l.m)
}
