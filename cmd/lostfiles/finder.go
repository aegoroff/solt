package lostfiles

import (
	"solt/cmd/api"
	"strings"
)

type finder struct {
	m api.Matcher
}

func newFinder(includedFiles, excludedFolders []string) (*finder, error) {
	excludes, err := api.NewPartialMatcher(excludedFolders, strings.ToUpper)
	if err != nil {
		return nil, err
	}

	includes := api.NewExactMatch(includedFiles)

	m := api.NewLostItemMatcher(includes, excludes)

	return &finder{m: m}, nil
}

func (l *finder) find(files []string) []string {
	return api.Filter(files, l.m)
}
