package fw

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NoMatch_False(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	data := generateRandomStringSlice(size, 50)
	bloom := NewBloomFilter(data)
	m := NewExactMatch(data)
	all := NewMatchAll(bloom, m)

	// Act
	ok := all.Match(nomatch)

	// Assert
	ass.False(ok)
}

func Test_Match_False(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	data := generateRandomStringSlice(size, 50)
	bloom := NewBloomFilter(data)
	m := NewExactMatch(data)
	all := NewMatchAll(bloom, m)
	s := data[size/2]

	// Act
	ok := all.Match(s)

	// Assert
	ass.True(ok)
}
