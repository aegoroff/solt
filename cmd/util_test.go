package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_sortAndOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items := []string{"rr", "aa", "xy"}
	mockp := newMockPrn()

	// Act
	sortAndOutput(mockp, items)

	// Assert
	ass.Equal(" aa\n rr\n xy\n", mockp.(*mockprn).String())
}

func Test_outputSortedMap(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items1 := []string{"rr", "gt", "xy"}
	items2 := []string{"ff", "lz", "xy"}

	m := map[string][]string{"a": items1, "b": items2}
	mockp := newMockPrn()

	// Act
	outputSortedMap(mockp, m, "SI")

	// Assert
	ass.Equal("\n<gray>SI: a</>\n gt\n rr\n xy\n\n<gray>SI: b</>\n ff\n lz\n xy\n", mockp.(*mockprn).String())
}
