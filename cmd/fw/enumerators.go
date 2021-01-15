package fw

import "solt/msvc"

// SolutionSlice provides *msvc.VisualStudioSolution slice
type SolutionSlice []*msvc.VisualStudioSolution

// Foreach enumerates all solutions within slice and
// calls each action on each
func (s SolutionSlice) Foreach(actions ...Solutioner) {
	for _, sol := range s {
		for _, h := range actions {
			h.Solution(sol)
		}
	}
}
