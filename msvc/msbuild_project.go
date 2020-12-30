package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
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
