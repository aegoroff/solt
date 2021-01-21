package in

import (
	"github.com/stretchr/testify/assert"
	"solt/solution"
	"testing"
)

func Test_configurationPlatform_allow(t *testing.T) {
	var tests = []struct {
		name   string
		expect bool
	}{
		{"SolutionConfigurationPlatforms", true},
		{"a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			s := newSectioner()
			sect := &solution.Section{
				Name: tt.name,
			}

			// Act
			result := s.allow(sect)

			// Assert
			ass.Equal(tt.expect, result)
		})
	}
}

func Test_configurationPlatform_handle(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	s := newSectioner()
	sect := &solution.Section{
		Items: []*solution.SectionItem{
			{Key: "Debug|Any CPU"},
			{Key: "Release|Any CPU"},
			{Key: "Debug|x86"},
			{Key: "Release|x86"},
		},
	}

	// Act
	s.run(sect)

	// Assert
	ass.ElementsMatch([]string{"Debug", "Release"}, s.configurations.Items())
	ass.ElementsMatch([]string{"Any CPU", "x86"}, s.platforms.Items())
}
