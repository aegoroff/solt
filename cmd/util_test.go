package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_sortAndOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items := []string{"rr", "aa", "xy"}
	buf := bytes.NewBufferString("")

	// Act
	sortAndOutput(buf, items)

	// Assert
	ass.Equal(" aa\n rr\n xy\n", buf.String())
}

func Test_outputSortedMap(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items1 := []string{"rr", "gt", "xy"}
	items2 := []string{"ff", "lz", "xy"}

	m := map[string][]string{"a": items1, "b": items2}
	buf := bytes.NewBufferString("")

	// Act
	outputSortedMap(buf, m, "SI")

	// Assert
	ass.Equal("\nSI: a\n gt\n rr\n xy\n\nSI: b\n ff\n lz\n xy\n", buf.String())
}
