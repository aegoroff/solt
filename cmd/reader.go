package cmd

import (
	"github.com/aegoroff/godatastruct/collections"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
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

type folderContent struct {
	packages  *Packages
	projects  []*msbuildProject
	solutions []*visualStudioSolution
}

type folder struct {
	content *folderContent
	path    string
}

func (x *folder) LessThan(y interface{}) bool {
	return sortfold.CompareFold(x.path, (y.(*folder)).path) < 0
}

func (x *folder) EqualTo(y interface{}) bool {
	return strings.EqualFold(x.path, (y.(*folder)).path)
}

func newTreeNode(f *folder) *rbtree.Comparable {
	var r rbtree.Comparable
	r = f
	return &r
}

func selectSolutions(foldersTree *rbtree.RbTree) []*visualStudioSolution {
	var solutions []*visualStudioSolution
	// Select only folders that contain solution(s)
	foldersTree.WalkInorder(func(n *rbtree.Node) {
		f := (*n.Key).(*folder)
		content := f.content
		if len(content.solutions) == 0 {
			return
		}
		for _, sln := range content.solutions {
			solutions = append(solutions, sln)
		}
	})
	return solutions
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

func readProjectDir(path string, fs afero.Fs, action func(we *walkEntry)) *rbtree.RbTree {
	result := rbtree.NewRbTree()

	aggregateChannel := make(chan *folder, 4)
	slowReadChannel := make(chan *walkEntry, 16)

	var wg sync.WaitGroup

	// Aggregating goroutine
	go func() {
		defer wg.Done()
		for f := range aggregateChannel {
			key := newTreeNode(f)
			if current, ok := result.Search(key); !ok {
				// Create new node
				n := rbtree.NewNode(key)
				result.Insert(n)
			} else {
				// Update folder node that has already been created before
				content := (*current.Key).(*folder).content

				if f.content.packages != nil {
					content.packages = f.content.packages
				} else {
					content.projects = append(content.projects, f.content.projects...)
					content.solutions = append(content.solutions, f.content.solutions...)
				}
			}
		}
	}()

	// Reading files goroutine
	go func() {
		defer close(aggregateChannel)

		for we := range slowReadChannel {
			if strings.EqualFold(we.Name, packagesConfigFile) {
				if folder, ok := onPackagesConfig(we, fs); ok {
					aggregateChannel <- folder
				}
			}

			ext := filepath.Ext(we.Name)
			if strings.EqualFold(ext, csharpProjectExt) || strings.EqualFold(ext, cppProjectExt) {
				if folder, ok := onMsbuildProject(we, fs); ok {
					aggregateChannel <- folder
				}
			}

			if strings.EqualFold(ext, solutionFileExt) {
				if folder, ok := onSolution(we, fs); ok {
					aggregateChannel <- folder
				}
			}
		}
	}()

	// Start reading path
	wg.Add(1)
	walkDirBreadthFirst(path, fs, func(parent string, entry os.FileInfo) {
		if entry.IsDir() {
			return
		}

		we := &walkEntry{IsDir: false, Size: entry.Size(), Parent: parent, Name: entry.Name()}
		slowReadChannel <- we

		action(we)
	})

	close(slowReadChannel)

	wg.Wait()

	return result
}

// Create packages model from packages.config
func onPackagesConfig(we *walkEntry, fs afero.Fs) (*folder, bool) {
	pack := Packages{}

	err := onXmlFile(we, fs, &pack)
	if err != nil {
		return nil, false
	}

	f := createFolder(we)

	f.content.packages = &pack

	return f, true
}

// Create project model from project file
func onMsbuildProject(we *walkEntry, fs afero.Fs) (*folder, bool) {
	project := Project{}

	err := onXmlFile(we, fs, &project)
	if err != nil {
		return nil, false
	}

	f := createFolder(we)

	p := msbuildProject{project: &project, path: filepath.Join(we.Parent, we.Name)}

	f.content.projects = append(f.content.projects, &p)

	return f, true
}

// Create solution model from file
func onSolution(we *walkEntry, fs afero.Fs) (*folder, bool) {
	solpath := filepath.Join(we.Parent, we.Name)
	reader, err := fs.Open(filepath.Clean(solpath))
	if err != nil {
		log.Println(err)
		return nil, false
	}

	sln, err := solution.Parse(reader)

	if err != nil {
		log.Println(err)
		return nil, false
	}

	f := createFolder(we)

	s := visualStudioSolution{solution: sln, path: filepath.Join(we.Parent, we.Name)}

	f.content.solutions = append(f.content.solutions, &s)

	return f, true
}

func createFolder(we *walkEntry) *folder {
	f := folder{
		content: &folderContent{
			solutions: []*visualStudioSolution{},
			projects:  []*msbuildProject{},
		},
		path: we.Parent,
	}
	return &f
}

func onXmlFile(we *walkEntry, fs afero.Fs, result interface{}) error {
	full := filepath.Join(we.Parent, we.Name)

	err := unmarshalXmlFrom(full, fs, result)
	if err != nil {
		log.Printf("%s: %v\n", full, err)
		return err
	}

	return nil
}
