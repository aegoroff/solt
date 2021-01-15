package info

import "solt/msvc"

type solutioner interface {
	solution(sl *msvc.VisualStudioSolution)
}
