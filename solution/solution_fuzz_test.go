package solution

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func FuzzParseSolution(f *testing.F) {
	f.Add(Vs2013)
	f.Add(Vs2010)
	f.Fuzz(func(t *testing.T, orig string) {
		// Arrange
		ass := assert.New(t)

		// Act
		sol := parse(orig, false)

		// Assert
		ass.NotNil(sol)
	})
}
