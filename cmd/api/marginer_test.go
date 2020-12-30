package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarginer_Margin(t *testing.T) {
	// Arrange
	const s = "x"
	ass := assert.New(t)
	var tests = []struct {
		value  int
		input  string
		result string
	}{
		{0, "", ""},
		{1, "", " "},
		{0, s, "x"},
		{1, s, " x"},
		{2, s, "  x"},
		{-1, s, "x"},
	}

	for _, test := range tests {
		t.Run(test.result, func(t *testing.T) {
			// Act
			m := NewMarginer(test.value)
			result := m.Margin(test.input)

			// Assert
			ass.Equal(test.result, result)
		})
	}
}