package msvc

import (
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"path/filepath"
	"solt/internal/sys"
	"strings"
)

// Folder defines filesystem folder descriptor (path and content structure)
type Folder struct {
	Content *FolderContent
	Path    string
}

// NugetPackage defines nu package descriptor
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

// NugetPackages gets all nu packages found in a folder
func (c *FolderContent) NugetPackages() ([]*NugetPackage, []string) {
	result := make([]*NugetPackage, 0)
	var sources []string

	if c.Packages != nil {
		// old style projects (nu packages references in separate files)
		result = append(result, c.Packages.nugetPackages()...)
		sources = append(sources, packagesConfigFile)
	}
	for _, prj := range c.Projects {
		np := prj.Project.nugetPackages()
		result = append(result, np...)
		if len(np) == 0 {
			continue
		}
		pp := sys.ToValidPath(prj.path)
		sources = append(sources, filepath.Base(pp))
	}
	return result, sources
}

// Less implements rbtree.Comparable interface
func (x *Folder) Less(y rbtree.Comparable) bool {
	return sortfold.CompareFold(x.Path, y.(*Folder).Path) < 0
}

// Equal implements rbtree.Comparable interface
func (x *Folder) Equal(y rbtree.Comparable) bool {
	return strings.EqualFold(x.Path, y.(*Folder).Path)
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
