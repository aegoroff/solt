package ux

import (
	"github.com/stretchr/testify/assert"
	"solt/internal/out"
	"testing"
)

func Test_writeSlice(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items := []string{"rr", "aa", "xy"}
	env := out.NewMemoryEnvironment()

	p := out.NewPrinter(env)
	s := NewScreener(p)

	// Act
	s.WriteSlice(items)

	// Assert
	ass.Equal("  aa\n  rr\n  xy\n", env.String())
}

func Test_writeMap(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items1 := []string{"rr", "gt", "xy"}
	items2 := []string{"ff", "lz", "xy"}

	m := map[string][]string{"a": items1, "b": items2}
	env := out.NewMemoryEnvironment()
	p := out.NewPrinter(env)
	s := NewScreener(p)

	// Act
	s.WriteMap(m, "SI")

	// Assert
	ass.Equal("\n SI: a\n  gt\n  rr\n  xy\n\n SI: b\n  ff\n  lz\n  xy\n", env.String())
}
