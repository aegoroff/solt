package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_SuccessMatchOneOfPatterns_ResultTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	patterns := []string{"xxx", "yyy", "zzz"}
	m := createAhoCorasickMachine(patterns)

	// Act
	result := Match(m, "yyyyy")

	// Assert
	ass.True(result)
}

func Test_SuccessMatchOneOfPatternsExactly_ResultTrue(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	patterns := []string{"xxx", "yyy", "zzz"}
	m := createAhoCorasickMachine(patterns)

	// Act
	result := Match(m, "yyy")

	// Assert
	ass.True(result)
}

func Test_NotMatchAnyOfPatterns_ResultFalse(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	patterns := []string{"xxx", "yyy", "zzz"}
	m := createAhoCorasickMachine(patterns)

	// Act
	result := Match(m, "aaa")

	// Assert
	ass.False(result)
}
