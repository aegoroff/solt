package info

import (
	"github.com/stretchr/testify/assert"
	"solt/solution"
	"testing"
)

func Test_sections_foreach(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	act := newSectioner()
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
