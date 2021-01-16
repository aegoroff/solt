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
	for _, sp := range s.Solution.Projects {
		if sp.TypeID != solution.IDSolutionFolder {
			fullProjectPath := filepath.Join(solutionPath, sp.Path)
			paths[i] = decorator(fullProjectPath)
			i++
		}
	}
	return paths[0:i]
}
