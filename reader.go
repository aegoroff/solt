package main

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func readProjectDir(path string, action func(we *walkEntry)) *rbtree.RbTree {
	readch := make(chan *walkEntry, 1024)

	go func(ch chan<- *walkEntry) {
		walkDirBreadthFirst(path, func(parent string, entry os.FileInfo) {
			if entry.IsDir() {
				return
			}

			ch <- &walkEntry{IsDir: false, Size: entry.Size(), Parent: parent, Name: entry.Name()}
		})
		close(ch)
	}(readch)

	result := rbtree.NewRbTree()

	aggregatech := make(chan *folder, 1024)

	go func(ch <-chan *folder) {
		for {
			f, ok := <-ch
			if !ok {
				break
			}
			key := createProjectTreeNode(f.path, f.info)

			if current, ok := rbtree.Search(result.Root, key); !ok {
				n := rbtree.NewNode(*key)
				rbtree.Insert(result, n)
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
	}(aggregatech)

	for {
		we, ok := <-readch
		if !ok {
			close(aggregatech)
			break
		}

		if we.Name == packagesConfigFile {
			if f, ok := onPackagesConfig(we); ok {
				aggregatech <- f
			}
		}

		ext := strings.ToLower(filepath.Ext(we.Name))
		if ext == csharpProjectExt || ext == cppProjectExt {
			if f, ok := onMsbuildProject(we); ok {
				aggregatech <- f
			}
		}

		action(we)
	}

	return result
}

func onPackagesConfig(we *walkEntry) (*folder, bool) {
	// Create packages model from packages.config
	pack := Packages{}

	f, ok := onXmlFile(we, &pack)
	if !ok {
		return nil, false
	}

	f.info.packages = &pack

	return f, true
}

func onMsbuildProject(we *walkEntry) (*folder, bool) {
	// Create project model from project file
	project := Project{}

	f, ok := onXmlFile(we, &project)
	if !ok {
		return nil, false
	}

	f.info.project = &project

	return f, true
}

func onXmlFile(we *walkEntry, result interface{}) (*folder, bool) {
	full := filepath.Join(we.Parent, we.Name)

	err := unmarshalXml(full, result)
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
