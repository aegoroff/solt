package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	var tests = []struct {
		cmd string
	}{
		{"version"},
		{"ver"},
	}

	for _, test := range tests {
		// Arrange
		ass := assert.New(t)
		memfs := afero.NewMemMapFs()
		p := newMockPrn()

		// Act
		_ = Execute(memfs, p.w, test.cmd)

		// Assert
		ass.Equal(fmt.Sprintf("solt v%s\n", Version), p.w.String())
	}
}
