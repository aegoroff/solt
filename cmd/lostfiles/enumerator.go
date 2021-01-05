package lostfiles

import (
	"solt/msvc"
)

type enumerator struct {
	skipper  *skpper
	includes []string
	exister  exister
}

func newEnumerator(foldersToIgnore *skpper, exister exister) *enumerator {
	return &enumerator{
		skipper: foldersToIgnore,
		exister: exister,
	}
}

func (e *enumerator) includedFiles() []string {
	return e.includes
}

func (e *enumerator) enumerate(projects []*msvc.MsbuildProject) {
	for _, prj := range projects {
		e.skipper.fromProject(prj)

		includes := prj.Files()
		e.includes = append(e.includes, includes...)

		e.exister.exist(prj.Path, includes)
	}
}
