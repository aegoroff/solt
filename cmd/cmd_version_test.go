package cmd

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			memfs := afero.NewMemMapFs()
			p := newMockPrn()

			// Act
			_ = Execute(memfs, p.w, test.cmd...)

			// Assert
			ass.Contains(p.w.String(), fmt.Sprintf("solt v%s\n", Version))
		})
	}
}
