package cmd

import (
	"bytes"
	"github.com/anknown/ahocorasick"
)

// Match do string matching to several patterns that defined by goahocorasick.Machine
func Match(m *goahocorasick.Machine, s string) bool {
	terms := m.MultiPatternSearch([]rune(s), true)
	return len(terms) > 0
}

func newAhoCorasickMachine(matches []string) (*goahocorasick.Machine, error) {
	var runes [][]rune
	for _, s := range matches {
		runes = append(runes, bytes.Runes([]byte(s)))
	}
	machine := new(goahocorasick.Machine)
	err := machine.Build(runes)
	if err != nil {
		return nil, err
	}
	return machine, nil
}
