package info

import (
	"solt/msvc"
	"solt/solution"
)

type projectGroupper struct {
	byType map[string]int
}

func newProjectGroupper() *projectGroupper {
	return &projectGroupper{}
}

func (p *projectGroupper) ByType() map[string]int {
	return p.byType
}

func (p *projectGroupper) Solution(vs *msvc.VisualStudioSolution) {
	p.byType = make(map[string]int)

	for _, pr := range vs.Solution.Projects {
		if pr.TypeID != solution.IDSolutionFolder {
			if pr.Type != "" {
				p.byType[pr.Type]++
			} else {
				p.byType[pr.TypeID]++
			}
		}
	}
}
