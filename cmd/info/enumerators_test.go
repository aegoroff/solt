package info

import (
	"github.com/stretchr/testify/assert"
	"solt/cmd/fw"
	"solt/msvc"
	"solt/solution"
	"sort"
	"testing"
)

func Test_sections_foreach(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	act := newConfigurationPlatform()
	sect := &solution.Section{
		Name: "SolutionConfigurationPlatforms",
		Items: []*solution.SectionItem{
			{Key: "Debug|Any CPU"},
		},
	}

	// Act
	sections{sect}.foreach(act)

	// Assert
	ass.ElementsMatch([]string{"Debug"}, act.configurations.Items())
	ass.ElementsMatch([]string{"Any CPU"}, act.platforms.Items())
}

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
			var slice fw.SolutionSlice
			for _, s := range tst.in {
				slice = append(slice, msvc.NewVisualStudioSolution(s))
			}

			// Act
			sort.Sort(slice)

			// Assert
			var result []string
			for _, s := range slice {
				result = append(result, s.Path())
			}
			ass.Equal(tst.expect, result)
		})
	}
}
