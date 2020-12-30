package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_writeSlice(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items := []string{"rr", "aa", "xy"}
	env := NewMemoryEnvironment()

	p := NewPrinter(env)
	s := NewScreener(p)

	// Act
	s.WriteSlice(items)

	// Assert
	ass.Equal(" aa\n rr\n xy\n", env.String())
}

func Test_writeMap(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items1 := []string{"rr", "gt", "xy"}
	items2 := []string{"ff", "lz", "xy"}

	m := map[string][]string{"a": items1, "b": items2}
	env := NewMemoryEnvironment()
	p := NewPrinter(env)
	s := NewScreener(p)

	// Act
	s.WriteMap(m, "SI")

	// Assert
	ass.Equal("\nSI: a\n gt\n rr\n xy\n\nSI: b\n ff\n lz\n xy\n", env.String())
}
