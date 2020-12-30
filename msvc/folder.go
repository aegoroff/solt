package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"path/filepath"
	"solt/solution"
)

// Folder defines filesystem folder descriptor (path and content structure)
type Folder struct {
	Content *FolderContent
	Path    string
}

// NugetPackage defines nuget package descriptor
type NugetPackage struct {
	ID      string
	Version string
}

// FolderContent defines a filesystem folder information about
// it's MSVC content (solutions, projects, etc.)
type FolderContent struct {
	Packages  *packages
	Projects  []*MsbuildProject
	Solutions []*VisualStudioSolution
}

// NugetPackages gets all nuget packages found in a folder
func (c *FolderContent) NugetPackages() ([]*NugetPackage, []string) {
	result := make([]*NugetPackage, 0)
	var sources []string

	if c.Packages != nil {
		// old style projects (nuget packages references in separate files)
		result = append(result, c.Packages.nugetPackages()...)
		sources = append(sources, packagesConfigFile)
	}
	for _, prj := range c.Projects {
		np := prj.Project.nugetPackages()
		result = append(result, np...)
		if len(np) == 0 {
			continue
		}
		pp := solution.ToValidPath(prj.Path)
		sources = append(sources, filepath.Base(pp))
	}
	return result, sources
}

// LessThan implements rbtree.Comparable interface
func (x *Folder) LessThan(y rbtree.Comparable) bool {
	return x.compare(y) < 0
}

// EqualTo implements rbtree.Comparable interface
func (x *Folder) EqualTo(y rbtree.Comparable) bool {
	return x.compare(y) == 0
}

func (x *Folder) compare(y rbtree.Comparable) int {
	return sortfold.CompareFold(x.Path, y.(*Folder).Path)
}

func newFolder(path string) *Folder {
	f := Folder{
		Content: &FolderContent{
			Solutions: []*VisualStudioSolution{},
			Projects:  []*MsbuildProject{},
		},
		Path: filepath.Dir(path),
	}
	return &f
}

func (x *Folder) copyContent(to *Folder) {
	toC := to.Content
	fromC := x.Content
	if fromC.Packages != nil {
		toC.Packages = fromC.Packages
	} else {
		toC.Projects = append(toC.Projects, fromC.Projects...)
		toC.Solutions = append(toC.Solutions, fromC.Solutions...)
	}
}