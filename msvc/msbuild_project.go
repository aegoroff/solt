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
	path    string
}

// NewMsbuildProject creates new *MsbuildProject instance
func NewMsbuildProject(path string) *MsbuildProject {
	return &MsbuildProject{path: path}
}

// Path gets full path to project
func (prj *MsbuildProject) Path() string {
	return prj.path
}

// IsSdkProject gets whether a project is a the new VS 2017 or later project
func (prj *MsbuildProject) IsSdkProject() bool {
	return prj.Project.isSdkProject()
}

// Less implements rbtree.Comparable interface
func (prj *MsbuildProject) Less(y rbtree.Comparable) bool {
	return sortfold.CompareFold(prj.path, y.(*MsbuildProject).path) < 0
}

// Equal implements rbtree.Comparable interface
func (prj *MsbuildProject) Equal(y rbtree.Comparable) bool {
	return strings.EqualFold(prj.path, y.(*MsbuildProject).path)
}

// Items gets all files included into MSBuild project
func (prj *MsbuildProject) Items() []string {
	folderPath := filepath.Dir(prj.path)

	msp := prj.Project

	// Sometimes ugly but rather fast
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

	// If SDK project nu packages included into project file
	for i, pkg := range p.PackageReferences {
		result[i] = &NugetPackage{ID: pkg.ID, Version: pkg.Version}
	}

	return result
}
