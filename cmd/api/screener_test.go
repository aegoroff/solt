package api

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_writeSlice(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items := []string{"rr", "aa", "xy"}
	w := bytes.NewBufferString("")
	env := NewStringEnvironment(w)

	p := NewPrinter(env)
	s := NewScreener(p)

	// Act
	s.WriteSlice(items)

	// Assert
	ass.Equal(" aa\n rr\n xy\n", w.String())
}

func Test_writeMap(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items1 := []string{"rr", "gt", "xy"}
	items2 := []string{"ff", "lz", "xy"}

	m := map[string][]string{"a": items1, "b": items2}
	w := bytes.NewBufferString("")
	env := NewStringEnvironment(w)
	p := NewPrinter(env)
	s := NewScreener(p)

	// Act
	s.WriteMap(m, "SI")

	// Assert
	ass.Equal("\nSI: a\n gt\n rr\n xy\n\nSI: b\n ff\n lz\n xy\n", w.String())
}
