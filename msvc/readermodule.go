package msvc

import (
	"github.com/aegoroff/dirstat/scan"
	"github.com/spf13/afero"
	"log"
	"path/filepath"
	"solt/internal/sys"
	"solt/solution"
	"strings"
)

func newReaderModules(fs afero.Fs) []readerModule {
	pack := readerPackagesConfig{fs}
	msbuild := readerMsbuild{fs}
	sol := readerSolution{fs}
	return []readerModule{&pack, &msbuild, &sol}
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

func (r *reader) Handler(path string) {
	for _, m := range r.modules {
		if !m.filter(path) {
			continue
		}
		if folder, ok := m.read(path); ok {
			r.aggregator <- folder
		}
	}
}

// packages.config

func (*readerPackagesConfig) filter(path string) bool {
	_, file := filepath.Split(path)
	return strings.EqualFold(file, packagesConfigFile)
}

func (r *readerPackagesConfig) read(path string) (*Folder, bool) {
	pack := packages{}
	d := sys.NewXMLDecoder(nil)

	err := d.UnmarshalFrom(path, r.fs, &pack)
	if err != nil {
		return nil, false
	}

	f := newFolder(path)

	f.Content.Packages = &pack

	return f, true
}

// MSBuild projects

func (*readerMsbuild) filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, csharpProjectExt) || strings.EqualFold(ext, cppProjectExt)
}

func (r *readerMsbuild) read(path string) (*Folder, bool) {
	project := msbuildProject{}
	d := sys.NewXMLDecoder(nil)

	err := d.UnmarshalFrom(path, r.fs, &project)
	if err != nil {
		return nil, false
	}

	f := newFolder(path)

	p := MsbuildProject{Project: &project, Path: path}

	f.Content.Projects = append(f.Content.Projects, &p)

	return f, true
}

// VS Solutions

func (*readerSolution) filter(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, SolutionFileExt)
}

func (r *readerSolution) read(path string) (*Folder, bool) {
	file, err := r.fs.Open(filepath.Clean(path))
	if err != nil {
		log.Println(err)
		return nil, false
	}
	defer scan.Close(file)

	sln, err := solution.Parse(file)

	if err != nil {
		log.Println(err)
		return nil, false
	}

	f := newFolder(path)

	s := VisualStudioSolution{Solution: sln, Path: path}

	f.Content.Solutions = append(f.Content.Solutions, &s)

	return f, true
}
