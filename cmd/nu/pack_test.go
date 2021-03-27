package nu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_match(t *testing.T) {
	var tests = []struct {
		name   string
		p1     *pack
		p2     *pack
		result bool
	}{
		{"single ver ids eq", newPack("p1", "1.0"), newPack("p1", "1.0"), true},
		{"many other ver ids eq", newPack("p1", "1.0"), newPack("p1", "1.0", "1.1"), true},
		{"many both ver ids eq", newPack("p1", "1.0", "1.2"), newPack("p1", "1.0", "1.1"), true},
		{"many both ver no intersection ids eq", newPack("p1", "2.0", "1.2"), newPack("p1", "1.0", "1.1"), false},
		{"versions match ids not eq", newPack("p1", "1.0", "1.2"), newPack("p2", "1.0", "1.1"), false},
		{"no other ver ids eq", newPack("p1", "1.0", "1.2"), newPack("p1"), false},
		{"no this ver ids eq", newPack("p1"), newPack("p1", "1.0", "1.2"), false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)

			// Act
			result := test.p1.match(test.p2)

			// Assert
			ass.Equal(test.result, result)
		})
	}
}
