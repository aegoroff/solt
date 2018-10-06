package main

import (
	"bytes"
	"github.com/anknown/ahocorasick"
)

// Match do string matching to several patterns that defined by goahocorasick.Machine
func Match(m *goahocorasick.Machine, s string) bool {
	terms := m.MultiPatternSearch([]rune(s), true)
	return len(terms) > 0
}

func createAhoCorasickMachine(matches []string) *goahocorasick.Machine {
	var runes [][]rune
	for _, s := range matches {
		runes = append(runes, bytes.Runes([]byte(s)))
	}
	machine := new(goahocorasick.Machine)
	machine.Build(runes)
	return machine
}
