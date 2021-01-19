package info

import (
	"solt/solution"
)

type sections []*solution.Section

func (s sections) foreach(action *sectioner) {
	for _, s := range s {
		if action.allow(s) {
			action.run(s)
		}
	}
}
