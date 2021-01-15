package info

import (
	"solt/msvc"
	"solt/solution"
)

type solutioner interface {
	solution(*msvc.VisualStudioSolution)
}

type sectioner interface {
	allow(*solution.Section) bool
	run(*solution.Section)
}
