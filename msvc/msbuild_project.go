package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"path/filepath"
)

// MsbuildProject defines MSBuild project structure
type MsbuildProject struct {
	Project *msbuildProject
	Path    string
}

// LessThan implements rbtree.Comparable interface
func (prj *MsbuildProject) LessThan(y rbtree.Comparable) bool {
	return prj.compare(y) < 0
}

// EqualTo implements rbtree.Comparable interface
func (prj *MsbuildProject) EqualTo(y rbtree.Comparable) bool {
	return prj.compare(y) == 0
}

func (prj *MsbuildProject) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(prj.Path, y.(*MsbuildProject).Path)
}

// Files gets all files included into MSBuild project
func (prj *MsbuildProject) Files() []string {
	folderPath := filepath.Dir(prj.Path)

	msp := prj.Project

	l := len(msp.Contents) + len(msp.Nones) + len(msp.CLCompiles) + len(msp.CLInclude) + len(msp.Compiles)
	includes := make([]include, 0, l)
	includes = append(includes, msp.Contents...)
	includes = append(includes, msp.Nones...)
	includes = append(includes, msp.CLCompiles...)
	includes = append(includes, msp.CLInclude...)
	includes = append(includes, msp.Compiles...)

	result := make([]string, len(includes))
	copy(result, createPathsFromIncludes(includes, folderPath))

	return result
}

func createPathsFromIncludes(paths []include, basePath string) []string {
	if paths == nil {
		return []string{}
	}

	result := make([]string, len(paths))

	for i, c := range paths {
		result[i] = filepath.Join(basePath, c.path())
	}

	return result
}

func (p *msbuildProject) nugetPackages() []*NugetPackage {
	if p.PackageReferences == nil {
		return []*NugetPackage{}
	}

	result := make([]*NugetPackage, len(p.PackageReferences))

	// If SDK project nuget packages included into project file
	for i, pkg := range p.PackageReferences {
		result[i] = &NugetPackage{ID: pkg.ID, Version: pkg.Version}
	}

	return result
}
