package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_MatchOneOfPatterns(t *testing.T) {
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
		m := createAhoCorasickMachine(test.patterns)
		result := Match(m, test.input)

		// Assert
		ass.Equal(test.result, result)
	}
}
