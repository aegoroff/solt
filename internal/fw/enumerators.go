package fw

import (
	"github.com/akutz/sortfold"
	"solt/msvc"
)

// SolutionSlice provides *msvc.VisualStudioSolution slice
type SolutionSlice []*msvc.VisualStudioSolution

func (s SolutionSlice) Len() int { return len(s) }

func (s SolutionSlice) Less(i, j int) bool {
	return sortfold.CompareFold(s[i].Path(), s[j].Path()) < 0
}

func (s SolutionSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Foreach enumerates all solutions within slice and
// calls each action on each
func (s SolutionSlice) Foreach(actions ...Solutioner) {
	for _, sol := range s {
		for _, h := range actions {
			h.Solution(sol)
		}
	}
}
