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

	sln.Projects(func(project *solution.Project) {
		p := msvc.NewMsbuildProject(filepath.Join(solutionPath, project.Path))

		found, ok := s.sdkProjects.Search(p)
		if ok {
			callFn(found.(*msvc.MsbuildProject))
		}
	})
}
