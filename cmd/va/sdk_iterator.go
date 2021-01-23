package va

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

type sdkIterator struct {
	sdkProjects rbtree.RbTree
}

func newSdkIterator(allProjects []*msvc.MsbuildProject) *sdkIterator {
	sdkProjects := rbtree.New()

	for _, p := range allProjects {
		if p.IsSdkProject() {
			sdkProjects.Insert(p)
		}
	}

	return &sdkIterator{
		sdkProjects: sdkProjects,
	}
}

func (s *sdkIterator) foreach(sln *msvc.VisualStudioSolution, callFn func(*msvc.MsbuildProject)) {
	solutionPath := filepath.Dir(sln.Path())
	for _, prj := range sln.Solution.Projects {
		if prj.TypeID == solution.IDSolutionFolder {
			continue
		}

		p := msvc.NewMsbuildProject(filepath.Join(solutionPath, prj.Path))

		found, ok := s.sdkProjects.Search(p)
		if !ok {
			continue
		}

		callFn(found.(*msvc.MsbuildProject))
	}
}
