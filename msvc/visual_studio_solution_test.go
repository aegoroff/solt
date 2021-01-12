package msvc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortSolutions(t *testing.T) {
	var tests = []struct {
		name   string
		in     []string
		expect []string
	}{
		{"several unsorted", []string{"b.sln", "a.sln"}, []string{"a.sln", "b.sln"}},
		{"one", []string{"a.sln"}, []string{"a.sln"}},
		{"none", []string{}, nil},
		{"several unsorted case test", []string{"ac.sln", "AB.sln"}, []string{"AB.sln", "ac.sln"}},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			var slice []*VisualStudioSolution
			for _, s := range tst.in {
				slice = append(slice, &VisualStudioSolution{path: s})
			}

			// Act
			SortSolutions(slice)

			// Assert
			var result []string
			for _, solution := range slice {
				result = append(result, solution.Path())
			}
			ass.Equal(tst.expect, result)
		})
	}
}
