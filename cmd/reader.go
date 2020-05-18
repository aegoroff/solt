package cmd

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/spf13/afero"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type folderInfo struct {
	packages    *Packages
	project     *Project
	projectPath *string
}

type folder struct {
	info *folderInfo
	path string
}

func getFilesIncludedIntoProject(info *folderInfo) []string {
	dir := filepath.Dir(*info.projectPath)
	var result []string
	result = append(result, getFiles(info.project.Contents, dir)...)
	result = append(result, getFiles(info.project.Nones, dir)...)
	result = append(result, getFiles(info.project.CLCompiles, dir)...)
	result = append(result, getFiles(info.project.CLInclude, dir)...)
	result = append(result, getFiles(info.project.Compiles, dir)...)
	return result
}

func getFiles(includes []Include, dir string) []string {
	if includes == nil {
		return []string{}
	}

	var result []string

	for _, c := range includes {
		fp := filepath.Join(dir, c.Path)
		result = append(result, fp)
	}

	return result
}

func readProjectDir(path string, fs afero.Fs, action func(we *walkEntry)) *rbtree.RbTree {
	result := rbtree.NewRbTree()

	aggregateChannel := make(chan *folder, 4)
	deferReadChannel := make(chan *walkEntry, 16)

	var wg sync.WaitGroup

	// Aggregating goroutine
	go func() {
		defer wg.Done()
		for f := range aggregateChannel {
			key := newProjectTreeNode(f.path, f.info)

			if current, ok := result.Search(key); !ok {
				n := rbtree.NewNode(key)
				result.Insert(n)
			} else {
				// Update folder node that has already been created before
				info := (*current.Key).(projectTreeNode).info
				if info.project == nil {
					// Project read after packages.config
					info.project = f.info.project
					info.projectPath = f.info.projectPath
				} else if info.packages == nil {
					// Project read before packages.config
					info.packages = f.info.packages
				}
			}
		}
	}()

	// Reading files goroutine
	go func() {
		defer close(aggregateChannel)

		for we := range deferReadChannel {
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
		}
	}()

	// Start reading path
	wg.Add(1)
	walkDirBreadthFirst(path, fs, func(parent string, entry os.FileInfo) {
		if entry.IsDir() {
			return
		}

		we := &walkEntry{IsDir: false, Size: entry.Size(), Parent: parent, Name: entry.Name()}
		deferReadChannel <- we

		action(we)
	})

	close(deferReadChannel)

	wg.Wait()

	return result
}

// Create packages model from packages.config
func onPackagesConfig(we *walkEntry, fs afero.Fs) (*folder, bool) {
	pack := Packages{}

	f, ok := onXmlFile(we, fs, &pack)
	if !ok {
		return nil, false
	}

	f.info.packages = &pack

	return f, true
}

// Create project model from project file
func onMsbuildProject(we *walkEntry, fs afero.Fs) (*folder, bool) {
	project := Project{}

	f, ok := onXmlFile(we, fs, &project)
	if !ok {
		return nil, false
	}

	f.info.project = &project

	return f, true
}

func onXmlFile(we *walkEntry, fs afero.Fs, result interface{}) (*folder, bool) {
	full := filepath.Join(we.Parent, we.Name)

	err := unmarshalXmlFrom(full, fs, result)
	if err != nil {
		log.Printf("%s: %v\n", full, err)
		return nil, false
	}

	f := folder{
		info: &folderInfo{projectPath: &full},
		path: we.Parent,
	}

	return &f, true
}
