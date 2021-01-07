package msvc

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/akutz/sortfold"
	"github.com/spf13/afero"
	"path/filepath"
	"solt/solution"
	"sort"
)

const (
	// SolutionFileExt defines visual studio extension
	SolutionFileExt    = ".sln"
	csharpProjectExt   = ".csproj"
	cppProjectExt      = ".vcxproj"
	packagesConfigFile = "packages.config"
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

// ProjectHandler defines project handler function prototype
type ProjectHandler func(*MsbuildProject, *Folder)

// StringDecorator defines string decorating function
type StringDecorator func(s string) string

type filesystem struct{ fs afero.Fs }

func newFs(fs afero.Fs) scan.Filesystem {
	return &filesystem{fs: fs}
}

func (f *filesystem) Open(path string) (scan.File, error) {
	return f.fs.Open(path)
}
