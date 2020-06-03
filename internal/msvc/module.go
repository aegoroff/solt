package msvc

import (
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"solt/internal/sys"
	"solt/solution"
	"strings"
)

type readerHandler interface {
	handler(path string)
}

// ReaderModule defines reader module interface
type ReaderModule interface {
	Filter(path string) bool
	Read(path string) (*Folder, bool)
}

type readerModules struct {
	modules    []ReaderModule
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
		if !m.Filter(path) {
			continue
		}
		if folder, ok := m.Read(path); ok {
			r.aggregator <- folder
		}
	}
}

// packages.config

func (r *readerPackagesConfig) Filter(path string) bool {
	_, file := filepath.Split(path)
	return strings.EqualFold(file, packagesConfigFile)
}

func (r *readerPackagesConfig) Read(path string) (*Folder, bool) {
	pack := packages{}

	err := onXMLFile(path, r.fs, &pack)
	if err != nil {
		return nil, false
	}

	f := createFolder(path)

	f.Content.Packages = &pack

	return f, true
}

// MSBuild projects

func (r *readerMsbuild) Filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, csharpProjectExt) || strings.EqualFold(ext, cppProjectExt)
}

func (r *readerMsbuild) Read(path string) (*Folder, bool) {
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

// VS Solutions

func (r *readerSolution) Filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, SolutionFileExt)
}

func (r *readerSolution) Read(path string) (*Folder, bool) {
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
