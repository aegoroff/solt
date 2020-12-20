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
		_ = afero.WriteFile(memfs, "a.txt", []byte("a"), 0644)
		f := NewFiler(memfs, bytes.NewBufferString(""))

		// Act
		res := f.CheckExistence(tst.in)

		// Assert
		ass.ElementsMatch(tst.expect, res)
	}
}

func TestFiler_Read(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	path := "a.txt"
	_ = afero.WriteFile(memfs, path, []byte("a"), 0644)
	f := NewFiler(memfs, bytes.NewBufferString(""))

	// Act
	buf := f.Read(path)

	// Assert
	ass.NotNil(buf)
	ass.Equal("a", string(buf.Bytes()))
}

func TestFiler_Read_NotExist(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	path := "a.txt"
	f := NewFiler(memfs, bytes.NewBufferString(""))

	// Act
	buf := f.Read(path)

	// Assert
	ass.Nil(buf)
}

func TestFiler_Write(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	f := NewFiler(memfs, bytes.NewBufferString(""))
	path := "a.txt"

	// Act
	f.Write(path, []byte("a"))

	// Assert
	content, err := afero.ReadFile(memfs, path)
	ass.NoError(err)
	ass.Equal("a", string(content))
}

func TestFiler_Write_Error(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	readOnly := afero.NewReadOnlyFs(memfs)
	f := NewFiler(readOnly, bytes.NewBufferString(""))
	path := "a.txt"

	// Act
	f.Write(path, []byte("a"))

	// Assert
	content, err := afero.ReadFile(memfs, path)
	ass.Error(err)
	ass.Nil(content)
}

func TestClose(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	f, err := memfs.Open("some")

	// Act
	Close(f)

	// Assert
	ass.Error(err)
}
