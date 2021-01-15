package fw

import (
	"github.com/akutz/sortfold"
	"solt/msvc"
)

// SolutionSlice provides *msvc.VisualStudioSolution slice
type SolutionSlice []*msvc.VisualStudioSolution

func (v SolutionSlice) Len() int { return len(v) }

func (v SolutionSlice) Less(i, j int) bool {
	return sortfold.CompareFold(v[i].Path(), v[j].Path()) < 0
}

func (v SolutionSlice) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// Foreach enumerates all solutions within slice and
// calls each action on each
func (s SolutionSlice) Foreach(actions ...Solutioner) {
	for _, sol := range s {
		for _, h := range actions {
			h.Solution(sol)
		}
	}
}
