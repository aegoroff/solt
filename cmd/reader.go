package cmd

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

type msbuildProject struct {
	project *Project
	path    string
}

type visualStudioSolution struct {
	solution *solution.Solution
	path     string
}

type readerHandler interface {
	handler(path string)
}

type readerModule interface {
	filter(path string) bool
	read(path string) (*folder, bool)
}

func selectAllSolutionProjectPaths(sln *visualStudioSolution, normalize bool) collections.StringHashSet {
	solutionPath := filepath.Dir(sln.path)
	var paths = make(collections.StringHashSet)
	for _, sp := range sln.solution.Projects {
		if sp.TypeId == solution.IdSolutionFolder {
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

func getFilesIncludedIntoProject(prj *msbuildProject) []string {
	var result []string
	folderPath := filepath.Dir(prj.path)
	result = append(result, createPaths(prj.project.Contents, folderPath)...)
	result = append(result, createPaths(prj.project.Nones, folderPath)...)
	result = append(result, createPaths(prj.project.CLCompiles, folderPath)...)
	result = append(result, createPaths(prj.project.CLInclude, folderPath)...)
	result = append(result, createPaths(prj.project.Compiles, folderPath)...)

	return result
}

func createPaths(paths []Include, basePath string) []string {
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

func readProjectDir(path string, fs afero.Fs, action func(we *walkEntry)) rbtree.RbTree {
	result := rbtree.NewRbTree()

	aggregateChannel := make(chan *folder, 4)
	slowReadChannel := make(chan *walkEntry, 16)

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
				content := current.Key().(*folder).content

				if f.content.packages != nil {
					content.packages = f.content.packages
				} else {
					content.projects = append(content.projects, f.content.projects...)
					content.solutions = append(content.solutions, f.content.solutions...)
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

		for we := range slowReadChannel {
			rh.handler(we.Path)
		}
	}(&rm)

	handlers := []sys.ScanHandler{func(evt *sys.ScanEvent) {
		if evt.File == nil {
			return
		}
		f := evt.File
		we := &walkEntry{Size: f.Size, Path: f.Path}
		slowReadChannel <- we
		action(we)
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
	aggregator chan *folder
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
func (r *readerPackagesConfig) read(path string) (*folder, bool) {
	pack := Packages{}

	err := onXmlFile(path, r.fs, &pack)
	if err != nil {
		return nil, false
	}

	f := createFolder(path)

	f.content.packages = &pack

	return f, true
}

func (r *readerMsbuild) filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, csharpProjectExt) || strings.EqualFold(ext, cppProjectExt)
}

// Create project model from project file
func (r *readerMsbuild) read(path string) (*folder, bool) {
	project := Project{}

	err := onXmlFile(path, r.fs, &project)
	if err != nil {
		return nil, false
	}

	f := createFolder(path)

	p := msbuildProject{project: &project, path: path}

	f.content.projects = append(f.content.projects, &p)

	return f, true
}

func (r *readerSolution) filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, solutionFileExt)
}

// Create solution model from file
func (r *readerSolution) read(path string) (*folder, bool) {
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

	s := visualStudioSolution{solution: sln, path: path}

	f.content.solutions = append(f.content.solutions, &s)

	return f, true
}

func createFolder(path string) *folder {
	f := folder{
		content: &folderContent{
			solutions: []*visualStudioSolution{},
			projects:  []*msbuildProject{},
		},
		path: filepath.Dir(path),
	}
	return &f
}

func onXmlFile(path string, fs afero.Fs, result interface{}) error {

	err := sys.UnmarshalXmlFrom(path, fs, result)
	if err != nil {
		log.Printf("%s: %v\n", path, err)
		return err
	}

	return nil
}
