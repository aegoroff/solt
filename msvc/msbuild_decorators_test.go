package msvc

import (
	"github.com/stretchr/testify/assert"
	"solt/internal/sys"
	"testing"
)

func Test_newMsbuildStandardPaths(t *testing.T) {
	var tests = []struct {
		path     string
		expected string
	}{
		{"b.txt", sys.ToValidPath("a\\b.txt")},
		{"", "a"},
		{"$(MSBuildProjectDirectory)b.txt", sys.ToValidPath("a\\b.txt")},
		{"$(MSBuildThisFileDirectory)b.txt", sys.ToValidPath("a\\b.txt")},
	}
	for _, tst := range tests {
		t.Run(tst.path, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			d := newMsbuildStandardPaths("a")

			// Act
			result := d.decorate(tst.path)

			// Assert
			ass.Equal(tst.expected, result)
		})
	}
}
