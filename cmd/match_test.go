package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MatchOneOfPatterns_Partial(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var tests = []struct {
		patterns []string
		input    string
		result   bool
	}{
		{[]string{"xxx", "yyy", "zzz"}, "yyyyy", true},
		{[]string{"xxx", "yyy", "zzz"}, "yyy", true},
		{[]string{"xxx", "yyy", "zzz"}, "yyyzzz", true},
		{[]string{"xxx", "yyy", "zzz"}, "cccyyybbb", true},
		{[]string{"xxx", "yyy", "zzz"}, "aaa", false},
	}

	for _, test := range tests {
		// Act
		m, _ := NewPartialMatcher(test.patterns)
		result := m.Match(test.input)

		// Assert
		ass.Equal(test.result, result)
	}
}

func Test_MatchOneOfPatterns_Exact(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var tests = []struct {
		patterns []string
		input    string
		result   bool
	}{
		{[]string{"xxx", "yyy", "zzz"}, "yyyyy", false},
		{[]string{"xxx", "yyy", "zzz"}, "yyy", true},
		{[]string{"xxx", "yyy", "zzz"}, "yyyzzz", false},
		{[]string{"xxx", "yyy", "zzz"}, "cccyyybbb", false},
		{[]string{"xxx", "yyy", "zzz"}, "aaa", false},
	}

	for _, test := range tests {
		// Act
		m := NewExactMatchS(test.patterns, func(s string) string { return s })
		result := m.Match(test.input)

		// Assert
		ass.Equal(test.result, result)
	}
}
