package sys

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckExistence(t *testing.T) {
	var tests = []struct {
		in     []string
		expect []string
	}{
		{[]string{"a.txt", "b.txt"}, []string{"b.txt"}},
		{[]string{"a.txt"}, []string{}},
	}
	for _, tst := range tests {
		// Arrange
		ass := assert.New(t)
		memfs := afero.NewMemMapFs()
		afero.WriteFile(memfs, "a.txt", []byte("a"), 0644)
		f := NewFiler(memfs, bytes.NewBufferString(""))

		// Act
		res := f.CheckExistence(tst.in)

		// Assert
		ass.ElementsMatch(tst.expect, res)
	}
}
