package msvc

import (
	"path/filepath"
	"solt/solution"
)

// VisualStudioSolution defines VS solution that contains *solution.Solution
// and it's path
type VisualStudioSolution struct {
	// Solution structure
	Solution *solution.Solution

	// filesystem path
	path string
}

// NewVisualStudioSolution creates new *VisualStudioSolution instance and assigns path to it
func NewVisualStudioSolution(path string) *VisualStudioSolution {
	return &VisualStudioSolution{path: path}
}

// Path gets full path to solution
func (s *VisualStudioSolution) Path() string {
	return s.path
}

// Items gets all paths of projects' included into solution
func (s *VisualStudioSolution) Items() []string {
	return s.AllProjectPaths(func(s string) string { return s })
}

// AllProjectPaths gets all possible projects' paths defined in solution
func (s *VisualStudioSolution) AllProjectPaths(decorator StringDecorator) []string {
	solutionPath := filepath.Dir(s.path)
	var paths = make([]string, len(s.Solution.Projects))
	i := 0
	s.Projects(func(project *solution.Project) {
		fullProjectPath := filepath.Join(solutionPath, project.Path)
		paths[i] = decorator(fullProjectPath)
		i++
	})
	return paths[:i]
}

// Projects enumerates all solutions projects that
// has type different from {2150E333-8FDC-42A3-9474-1A3956D46DE8} (solution folder)
func (s *VisualStudioSolution) Projects(callFn func(*solution.Project)) {
	for _, sp := range s.Solution.Projects {
		if sp.TypeID != solution.IDSolutionFolder {
			callFn(sp)
		}
	}
}
