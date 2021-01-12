package msvc

import (
	"github.com/akutz/sortfold"
	"path/filepath"
	"solt/solution"
	"sort"
)

// VisualStudioSolution defines VS solution that contains *solution.Solution
// and it's path
type VisualStudioSolution struct {
	// Solution structure
	Solution *solution.Solution

	// filesystem path
	path string
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
	var paths = make([]string, 0, len(s.Solution.Projects))
	for _, sp := range s.Solution.Projects {
		if sp.TypeID == solution.IDSolutionFolder {
			continue
		}
		fullProjectPath := filepath.Join(solutionPath, sp.Path)
		paths = append(paths, decorator(fullProjectPath))
	}
	return paths
}

// SortSolutions sorts solutions by path
func SortSolutions(solutions []*VisualStudioSolution) {
	sort.Sort(visualStudioSolutionSlice(solutions))
}

type visualStudioSolutionSlice []*VisualStudioSolution

func (v visualStudioSolutionSlice) Len() int { return len(v) }

func (v visualStudioSolutionSlice) Less(i, j int) bool {
	return sortfold.CompareFold(v[i].path, v[j].path) < 0
}

func (v visualStudioSolutionSlice) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
