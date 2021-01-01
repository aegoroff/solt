package msvc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_includes_nil(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	var i includes

	// Act
	r := i.paths("")

	// Assert
	ass.Empty(r)
}
