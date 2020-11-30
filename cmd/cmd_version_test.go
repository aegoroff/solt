package cmd

import (
	"fmt"
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

		appPrinter = newMockPrn()

		// Act
		_ = Execute(test.cmd)

		// Assert
		ass.Equal(fmt.Sprintf("solt v%s\n", Version), appPrinter.(*mockprn).String())
	}
}
