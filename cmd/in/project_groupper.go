package in

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

	vs.Projects(func(pr *solution.Project) {
		if pr.Type != "" {
			p.byType[pr.Type]++
		} else {
			p.byType[pr.TypeID]++
		}
	})
}
