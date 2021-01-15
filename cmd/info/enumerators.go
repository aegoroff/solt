package info

import (
	"solt/msvc"
	"solt/solution"
)

type solutions []*msvc.VisualStudioSolution
type sections []*solution.Section

func (s sections) foreach(action sectioner) {
	for _, s := range s {
		if action.allow(s) {
			action.run(s)
		}
	}
}

func (s solutions) foreach(actions []solutioner) {
	for _, sol := range s {
		for _, h := range actions {
			h.solution(sol)
		}
	}
}
