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
	incl := make(includes, l)
	copy(incl, msp.Contents)
	copy(incl[len(msp.Contents):], msp.Nones)
	copy(incl[len(msp.Contents)+len(msp.Nones):], msp.CLCompiles)
	copy(incl[len(msp.Contents)+len(msp.Nones)+len(msp.CLCompiles):], msp.CLInclude)
	copy(incl[len(msp.Contents)+len(msp.Nones)+len(msp.CLCompiles)+len(msp.CLInclude):], msp.Compiles)

	result := make([]string, len(incl))
	copy(result, incl.paths(folderPath))

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
