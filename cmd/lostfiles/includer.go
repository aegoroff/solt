package lostfiles

import "solt/msvc"

type includer struct {
	includes []string
	exister  exister
}

func newIncluder(exister exister) *includer {
	return &includer{exister: exister}
}

func (i *includer) fromProject(p *msvc.MsbuildProject) {
	includes := p.Files()
	i.includes = append(i.includes, includes...)

	i.exister.exist(p.Path, includes)
}

func (i *includer) files() []string {
	return i.includes
}
