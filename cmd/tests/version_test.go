package tests

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/cmd"
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

			// Act
			_ = cmd.Execute(memfs, w, test.cmd...)

			// Assert
			ass.Contains(w.String(), cmd.Version)
		})
	}
}

func Test_Help(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	w := bytes.NewBufferString("")

	// Act
	_ = cmd.Execute(memfs, w, "")

	// Assert
	ass.Contains(w.String(), "")

}
