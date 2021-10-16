package sys

import (
	"bytes"
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckExistence(t *testing.T) {
	var tests = []struct {
		name   string
		in     []string
		expect []string
	}{
		{"has unexist files", []string{"a.txt", "b.txt"}, []string{"b.txt"}},
		{"has unexist files in same dir", []string{"a.txt", "/b/b.txt", "/b/c.txt"}, []string{"/b/b.txt", "/b/c.txt"}},
		{"all files exist", []string{"a.txt"}, []string{}},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			memfs := afero.NewMemMapFs()
			_ = afero.WriteFile(memfs, "a.txt", []byte("a"), 0644)
			f := NewFiler(memfs, bytes.NewBufferString(""))

			// Act
			res := f.CheckExistence(tst.in)

			// Assert
			ass.ElementsMatch(tst.expect, res)
		})
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
	buf, err := f.Read(path)

	// Assert
	ass.NoError(err)
	ass.NotNil(buf)
	ass.Equal("a", string(buf))
}

func TestFiler_Read_NotExist(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	path := "a.txt"
	f := NewFiler(memfs, bytes.NewBufferString(""))

	// Act
	buf, err := f.Read(path)

	// Assert
	ass.Error(err)
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
	var tests = []struct {
		name    string
		content []byte
	}{
		{"read only fs", []byte("a")},
		{"nil", nil},
	}
	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			memfs := afero.NewMemMapFs()
			readOnly := afero.NewReadOnlyFs(memfs)
			f := NewFiler(readOnly, bytes.NewBufferString(""))
			path := "a.txt"

			// Act
			f.Write(path, tst.content)

			// Assert
			content, err := afero.ReadFile(memfs, path)
			ass.Error(err)
			ass.Nil(content)
		})
	}
}

func TestClose(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	f, err := memfs.Open("some")

	// Act
	scan.Close(f)

	// Assert
	ass.Error(err)
}

func TestFiler_Remove(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	_ = afero.WriteFile(memfs, "a.txt", []byte("a"), 0644)
	f := NewFiler(memfs, bytes.NewBufferString(""))
	path := "a.txt"

	// Act
	f.Remove([]string{path})

	// Assert
	_, err := afero.ReadFile(memfs, path)
	ass.Error(err)
}
