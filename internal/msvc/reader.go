package msvc

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"path/filepath"
	"solt/internal/sys"
	"solt/solution"
	"sync"
)

// SelectAllSolutionProjectPaths gets all possible projects' paths defined in solution
func SelectAllSolutionProjectPaths(sln *VisualStudioSolution, pathDecorator StringDecorator) collections.StringHashSet {
	solutionPath := filepath.Dir(sln.Path)
	var paths = make(collections.StringHashSet)
	for _, sp := range sln.Solution.Projects {
		if sp.TypeID == solution.IDSolutionFolder {
			continue
		}
		fullProjectPath := filepath.Join(solutionPath, sp.Path)
		paths.Add(pathDecorator(fullProjectPath))
	}
	return paths
}

// GetFilesIncludedIntoProject gets all files included into MSBuild project
func GetFilesIncludedIntoProject(prj *MsbuildProject) []string {
	var result []string
	folderPath := filepath.Dir(prj.Path)

	msp := prj.Project

	var includes []include
	includes = append(includes, msp.Contents...)
	includes = append(includes, msp.Nones...)
	includes = append(includes, msp.CLCompiles...)
	includes = append(includes, msp.CLInclude...)
	includes = append(includes, msp.Compiles...)

	result = append(result, createPathsFromIncludes(includes, folderPath)...)

	return result
}

// WalkProjects traverse all projects found in solution(s) folder
func WalkProjects(foldersTree rbtree.RbTree, action ProjectHandler) {
	w := &walkPrj{handler: action}
	walk(foldersTree, w)
}

// SelectSolutions gets all Visual Studion solutions found in a directory
func SelectSolutions(foldersTree rbtree.RbTree) []*VisualStudioSolution {
	w := walkSol{solutions: make([]*VisualStudioSolution, 0)}
	walk(foldersTree, &w)
	return w.solutions
}

// IsSdkProject gets whether a project is a the new VS 2017 or later project
func (p *msbuildProject) IsSdkProject() bool {
	if len(p.Sdk) > 0 {
		return true
	}
	if len(p.Imports) == 0 {
		return false
	}
	for _, imp := range p.Imports {
		if len(imp.Sdk) > 0 {
			return true
		}
	}
	return false
}

// ReadSolutionDir reads filesystem directory and all its childs to get information
// about all solutions and projects in this tree.
// It returns tree
func ReadSolutionDir(path string, fs afero.Fs, fileHandlers ...ReaderHandler) rbtree.RbTree {
	result := rbtree.NewRbTree()

	aggregateChannel := make(chan *Folder, 4)
	fileChannel := make(chan string, 16)

	var wg sync.WaitGroup

	// Aggregating goroutine
	go func() {
		defer wg.Done()
		for f := range aggregateChannel {
			if current, ok := result.Search(f); !ok {
				// Create new node
				result.Insert(f)
			} else {
				// Update folder node that has already been created before
				current := current.Key().(*Folder)
				merge(current, f)
			}
		}
	}()

	modules := newReaderModules(fs)

	rdr := reader{aggregator: aggregateChannel, modules: modules}

	fhandlers := []ReaderHandler{&rdr}
	fhandlers = append(fhandlers, fileHandlers...)

	// Reading files goroutine
	go func(handlers []ReaderHandler) {
		defer close(aggregateChannel)

		for path := range fileChannel {
			for _, h := range handlers {
				h.Handler(path)
			}
		}
	}(fhandlers)

	handlers := []sys.ScanHandler{func(evt *sys.ScanEvent) {
		if evt.File == nil {
			return
		}
		f := evt.File
		fileChannel <- f.Path
	}}

	// Start reading path
	wg.Add(1)

	sys.Scan(path, fs, handlers)

	close(fileChannel)

	wg.Wait()

	return result
}

func merge(to *Folder, from *Folder) {
	toC := to.Content
	fromC := from.Content
	if fromC.Packages != nil {
		toC.Packages = fromC.Packages
	} else {
		toC.Projects = append(toC.Projects, fromC.Projects...)
		toC.Solutions = append(toC.Solutions, fromC.Solutions...)
	}
}

func createPathsFromIncludes(paths []include, basePath string) []string {
	if paths == nil {
		return []string{}
	}

	var result []string

	for _, c := range paths {
		fp := filepath.Join(basePath, c.Path)
		result = append(result, fp)
	}

	return result
}