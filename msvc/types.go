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
	Path string
}

// SortSolutions sorts solutions by path
func SortSolutions(solutions []*VisualStudioSolution) {
	sort.Sort(visualStudioSolutionSlice(solutions))
}

type visualStudioSolutionSlice []*VisualStudioSolution

func (v visualStudioSolutionSlice) Len() int { return len(v) }

func (v visualStudioSolutionSlice) Less(i, j int) bool {
	return sortfold.CompareFold(v[i].Path, v[j].Path) < 0
}

func (v visualStudioSolutionSlice) Swap(i, j int) { v[i], v[j] = v[j], v[i] }

// AllProjectPaths gets all possible projects' paths defined in solution
func (s *VisualStudioSolution) AllProjectPaths(decorator StringDecorator) []string {
	solutionPath := filepath.Dir(s.Path)
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

// PassThrough just return original string without modification
// it can be used as non destructive decorator
func PassThrough(s string) string { return s }

type filesystem struct{ fs afero.Fs }

func newFs(fs afero.Fs) scan.Filesystem {
	return &filesystem{fs: fs}
}

func (f *filesystem) Open(path string) (scan.File, error) {
	return f.fs.Open(path)
}
