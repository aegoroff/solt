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

// packages.config

func (*readerPackagesConfig) allow(path string) bool {
	_, file := filepath.Split(path)
	return strings.EqualFold(file, packagesConfigFile)
}

func (r *readerPackagesConfig) read(path string, ch chan<- *Folder) {
	pack := packages{}
	d := sys.NewXMLDecoder(nil)

	err := d.UnmarshalFrom(path, r.fs, &pack)
	if err != nil {
		return
	}

	f := newFolder(path)

	f.Content.Packages = &pack

	ch <- f
}

// MSBuild projects

func (*readerMsbuild) allow(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, csharpProjectExt) || strings.EqualFold(ext, cppProjectExt)
}

func (r *readerMsbuild) read(path string, ch chan<- *Folder) {
	project := msbuildProject{}
	d := sys.NewXMLDecoder(nil)

	err := d.UnmarshalFrom(path, r.fs, &project)
	if err != nil {
		return
	}

	f := newFolder(path)

	p := MsbuildProject{Project: &project, path: path}

	f.Content.Projects = append(f.Content.Projects, &p)

	ch <- f
}

// VS Solutions

func (*readerSolution) allow(path string) bool {
	ext := filepath.Ext(path)
	return strings.EqualFold(ext, SolutionFileExt)
}

func (r *readerSolution) read(path string, ch chan<- *Folder) {
	file, err := r.fs.Open(filepath.Clean(path))
	if err != nil {
		log.Println(err)
		return
	}
	defer scan.Close(file)

	sln, err := solution.Parse(file)

	if err != nil {
		log.Println(err)
		return
	}

	f := newFolder(path)

	vs := NewVisualStudioSolution(path)
	vs.Solution = sln

	f.Content.Solutions = append(f.Content.Solutions, vs)

	ch <- f
}
