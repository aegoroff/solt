package msvc

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"solt/internal/sys"
	"solt/solution"
	"strings"
	"sync"
)

// MsbuildProject defines MSBuild project structure
type MsbuildProject struct {
	Project *msbuildProject
	Path    string
}

// VisualStudioSolution defines VS solution that contains *solution.Solution
// and it's path
type VisualStudioSolution struct {
	// Solution structure
	Solution *solution.Solution

	// filesystem pa
	Path string
}

type readerHandler interface {
	handler(path string)
}

type readerModule interface {
	filter(path string) bool
	read(path string) (*Folder, bool)
}

// SelectAllSolutionProjectPaths gets all possible projects' paths defined in solution
func SelectAllSolutionProjectPaths(sln *VisualStudioSolution, normalize bool) collections.StringHashSet {
	solutionPath := filepath.Dir(sln.Path)
	var paths = make(collections.StringHashSet)
	for _, sp := range sln.Solution.Projects {
		if sp.TypeID == solution.IDSolutionFolder {
			continue
		}
		fullProjectPath := filepath.Join(solutionPath, sp.Path)

		if normalize {
			key := strings.ToUpper(fullProjectPath)
			paths.Add(key)
		} else {
			paths.Add(fullProjectPath)
		}
	}
	return paths
}

// GetFilesIncludedIntoProject gets all files included into MSBuild project
func GetFilesIncludedIntoProject(prj *MsbuildProject) []string {
	var result []string
	folderPath := filepath.Dir(prj.Path)
	result = append(result, createPaths(prj.Project.Contents, folderPath)...)
	result = append(result, createPaths(prj.Project.Nones, folderPath)...)
	result = append(result, createPaths(prj.Project.CLCompiles, folderPath)...)
	result = append(result, createPaths(prj.Project.CLInclude, folderPath)...)
	result = append(result, createPaths(prj.Project.Compiles, folderPath)...)

	return result
}

func createPaths(paths []include, basePath string) []string {
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

// ReadSolutionDir reads filesystem directory and all its childs to get information
// about all solutions and projects in this tree.
// It returns tree
func ReadSolutionDir(path string, fs afero.Fs, action func(path string)) rbtree.RbTree {
	result := rbtree.NewRbTree()

	aggregateChannel := make(chan *Folder, 4)
	slowReadChannel := make(chan string, 16)

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
				content := current.Key().(*Folder).Content

				if f.Content.Packages != nil {
					content.Packages = f.Content.Packages
				} else {
					content.Projects = append(content.Projects, f.Content.Projects...)
					content.Solutions = append(content.Solutions, f.Content.Solutions...)
				}
			}
		}
	}()

	var modules []readerModule

	pack := readerPackagesConfig{fs}
	msbuild := readerMsbuild{fs}
	sol := readerSolution{fs}
	modules = append(modules, &pack, &msbuild, &sol)

	rm := readerModules{aggregator: aggregateChannel, modules: modules}

	// Reading files goroutine
	go func(rh readerHandler) {
		defer close(aggregateChannel)

		for path := range slowReadChannel {
			rh.handler(path)
		}
	}(&rm)

	handlers := []sys.ScanHandler{func(evt *sys.ScanEvent) {
		if evt.File == nil {
			return
		}
		f := evt.File
		slowReadChannel <- f.Path
		action(f.Path)
	}}

	// Start reading path
	wg.Add(1)

	sys.Scan(path, fs, handlers)

	close(slowReadChannel)

	wg.Wait()

	return result
}

type readerModules struct {
	modules    []readerModule
	aggregator chan *Folder
}

type readerPackagesConfig struct {
	fs afero.Fs
}

type readerMsbuild struct {
	fs afero.Fs
}

type readerSolution struct {
	fs afero.Fs
}

func (r *readerModules) handler(path string) {
	for _, m := range r.modules {
		if !m.filter(path) {
			continue
		}
		if folder, ok := m.read(path); ok {
			r.aggregator <- folder
		}
	}
}

func (r *readerPackagesConfig) filter(path string) bool {
	_, file := filepath.Split(path)
	return strings.EqualFold(file, packagesConfigFile)
}

// Create packages model from packages.config
func (r *readerPackagesConfig) read(path string) (*Folder, bool) {
	pack := packages{}

	err := onXMLFile(path, r.fs, &pack)
	if err != nil {
		return nil, false
	}

	f := createFolder(path)

	f.Content.Packages = &pack

	return f, true
}

func (r *readerMsbuild) filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, csharpProjectExt) || strings.EqualFold(ext, cppProjectExt)
}

// Create project model from project file
func (r *readerMsbuild) read(path string) (*Folder, bool) {
	project := msbuildProject{}

	err := onXMLFile(path, r.fs, &project)
	if err != nil {
		return nil, false
	}

	f := createFolder(path)

	p := MsbuildProject{Project: &project, Path: path}

	f.Content.Projects = append(f.Content.Projects, &p)

	return f, true
}

func (r *readerSolution) filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, solutionFileExt)
}

// Create solution model from file
func (r *readerSolution) read(path string) (*Folder, bool) {
	reader, err := r.fs.Open(filepath.Clean(path))
	if err != nil {
		log.Println(err)
		return nil, false
	}

	sln, err := solution.Parse(reader)

	if err != nil {
		log.Println(err)
		return nil, false
	}

	f := createFolder(path)

	s := VisualStudioSolution{Solution: sln, Path: path}

	f.Content.Solutions = append(f.Content.Solutions, &s)

	return f, true
}

func createFolder(path string) *Folder {
	f := Folder{
		Content: &FolderContent{
			Solutions: []*VisualStudioSolution{},
			Projects:  []*MsbuildProject{},
		},
		Path: filepath.Dir(path),
	}
	return &f
}

func onXMLFile(path string, fs afero.Fs, result interface{}) error {

	err := sys.UnmarshalXMLFrom(path, fs, result)
	if err != nil {
		log.Printf("%s: %v\n", path, err)
		return err
	}

	return nil
}
