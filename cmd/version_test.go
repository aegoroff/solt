package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/cmd/api"
	"testing"
)

func Test_Version(t *testing.T) {
	var tests = []struct {
		name string
		cmd  []string
	}{
		{"version", []string{"version"}},
		{"ver", []string{"ver"}},
		{"ver -d", []string{"ver", "-d"}},
		{"ver -d --cpuprofile", []string{"ver", "-d", "--cpuprofile", "cpu"}},
		{"ver -d --memprofile", []string{"ver", "-d", "--memprofile", "mem"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			memfs := afero.NewMemMapFs()
			w := bytes.NewBufferString("")
			env := api.NewStringEnvironment(w)

			// Act
			_ = Execute(memfs, env, test.cmd...)

			// Assert
			ass.Contains(w.String(), Version)
		})
	}
}

func Test_Help(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "")

	// Assert
	ass.Contains(w.String(), "")
}

func Test_Write_File(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)
	fp := "/f"

	// Act
	_ = Execute(memfs, env, "ver", "-f", fp)

	// Assert
	ass.Equal("", w.String())
	_, err := memfs.Stat(fp)
	ass.NoError(err)
}

func Test_Write_ReadonlyFile(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)
	fp := "/f"
	ro := afero.NewReadOnlyFs(memfs)

	// Act
	_ = Execute(ro, env, "ver", "-f", fp)

	// Assert
	ass.Contains(w.String(), Version)
	_, err := memfs.Stat(fp)
	ass.Error(err)
}
