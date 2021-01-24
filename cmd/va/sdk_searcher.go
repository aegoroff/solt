package va

import (
	"path/filepath"
	"solt/msvc"
	"solt/solution"
)

type sdkSearcher struct {
	projects     *sdkProjects
	solutionPath string
}

func newSdkSearcher(projects *sdkProjects, sln *msvc.VisualStudioSolution) *sdkSearcher {
	return &sdkSearcher{
		projects:     projects,
		solutionPath: filepath.Dir(sln.Path()),
	}
}

func (s *sdkSearcher) Search(prj *solution.Project) (*msvc.MsbuildProject, bool) {
	p := msvc.NewMsbuildProject(filepath.Join(s.solutionPath, prj.Path))
	return s.projects.search(p)
}
