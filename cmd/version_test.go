package cmd

import (
	"bytes"
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

		buf := bytes.NewBufferString("")
		appWriter = buf

		// Act
		rootCmd.SetArgs([]string{test.cmd})
		rootCmd.Execute()

		// Assert
		ass.Equal(fmt.Sprintf("solt v%s\n", Version), buf.String())
	}
}
