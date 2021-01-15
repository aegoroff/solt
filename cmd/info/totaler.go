package info

import (
	"solt/msvc"
)

type totaler struct {
	result *totals
}

func newTotaler() *totaler {
	return &totaler{
		result: &totals{},
	}
}

func (t *totaler) Solution(sl *msvc.VisualStudioSolution) {
	t.result.solutions++
	t.result.projects += len(sl.Solution.Projects)
}
