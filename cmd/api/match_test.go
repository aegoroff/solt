package api

import (
	"github.com/stretchr/testify/assert"
	"strings"
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
		{[]string{"xxx", "yyy", "zzz"}, "YYYYY", true},
		{[]string{"xxx", "yyy", "zzz"}, "yyy", true},
		{[]string{"xxx", "yyy", "zzz"}, "yyyzzz", true},
		{[]string{"xxx", "yyy", "zzz"}, "cccyyybbb", true},
		{[]string{"xxx", "yyy", "zzz"}, "CCCYYYBBB", true},
		{[]string{"xxx", "yyy", "zzz"}, "aaa", false},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			// Act
			m, _ := NewPartialMatcher(test.patterns, strings.ToUpper)
			result := m.Match(test.input)

			// Assert
			ass.Equal(test.result, result)
		})
	}
}

func TestNewPartialMatcher_EmptyMatches(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	// Act
	m, err := NewPartialMatcher([]string{}, strings.ToUpper)

	// Assert
	ass.Error(err)
	ass.Nil(m)
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
		t.Run(test.input, func(t *testing.T) {
			// Act
			m := NewExactMatch(test.patterns)
			result := m.Match(test.input)

			// Assert
			ass.Equal(test.result, result)
		})
	}
}

func Test_MatchAnyOneOfPatterns_Exact(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var tests = []struct {
		name     string
		patterns []string
		input    []string
		result   bool
	}{
		{"one matched", []string{"xxx", "yyy", "zzz"}, []string{"yyy"}, true},
		{"all matched", []string{"xxx", "yyy", "zzz"}, []string{"xxx", "yyy", "zzz"}, true},
		{"all not matched", []string{"eee", "rr", "dd"}, []string{"xxx", "yyy", "zzz"}, false},
		{"one not matched", []string{"xxx", "yyy", "zzz"}, []string{"eee"}, false},
		{"empty", []string{"xxx", "yyy", "zzz"}, []string{}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Act
			m := NewExactMatch(test.patterns)
			result := MatchAny(m, test.input)

			// Assert
			ass.Equal(test.result, result)
		})
	}
}
