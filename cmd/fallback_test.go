package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fallback(t *testing.T) {
	var tests = []struct {
		in     string
		expect int
	}{
		{"", 0},
		{"<x y=\"z\"/>", 0},
		{"\n<x y=\"z\"/>", 0},
		{"\n\n<x y=\"z\"/>", 1},
		{">\n\n<x y=\"z\"/>", 2},
		{"\r\n<x y=\"z\"/>", 0},
		{">\r\n<x y=\"z\"/>", 1},
		{"><x y=\"z\"/>", 1},
		{"    <x y=\"z\"/>", 0},
	}
	for _, tst := range tests {
		t.Run(tst.in, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)

			// Act
			res := fallback([]byte(tst.in))

			// Assert
			ass.Equal(tst.expect, res)
		})
	}
}
