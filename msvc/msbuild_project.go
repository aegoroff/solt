package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"path/filepath"
	"strings"
)

// MsbuildProject defines MSBuild project structure
type MsbuildProject struct {
	Project *msbuildProject
	Path    string
}

// LessThan implements rbtree.Comparable interface
func (prj *MsbuildProject) LessThan(y rbtree.Comparable) bool {
	return sortfold.CompareFold(prj.Path, y.(*MsbuildProject).Path) < 0
}

// EqualTo implements rbtree.Comparable interface
func (prj *MsbuildProject) EqualTo(y rbtree.Comparable) bool {
	return strings.EqualFold(prj.Path, y.(*MsbuildProject).Path)
}

// Files gets all files included into MSBuild project
func (prj *MsbuildProject) Files() []string {
	folderPath := filepath.Dir(prj.Path)

	msp := prj.Project

	l := len(msp.Contents) + len(msp.Nones) + len(msp.CLCompiles) + len(msp.CLInclude) + len(msp.Compiles)
	incl := make(includes, l)
	start := 0
	start += incl.append(msp.Contents, start)
	start += incl.append(msp.Nones, start)
	start += incl.append(msp.CLCompiles, start)
	start += incl.append(msp.CLInclude, start)
	incl.append(msp.Compiles, start)

	result := make([]string, l)
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
