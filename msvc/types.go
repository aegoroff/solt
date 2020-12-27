package msvc

import (
	"encoding/xml"
	"github.com/aegoroff/dirstat/scan"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"github.com/spf13/afero"
	"path/filepath"
	"solt/solution"
)

const (
	// SolutionFileExt defines visual studio extension
	SolutionFileExt    = ".sln"
	csharpProjectExt   = ".csproj"
	cppProjectExt      = ".vcxproj"
	packagesConfigFile = "packages.config"
)

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

func (p *packages) nugetPackages() []*NugetPackage {
	result := make([]*NugetPackage, len(p.Packages))

	for i, pkg := range p.Packages {
		result[i] = &NugetPackage{ID: pkg.ID, Version: pkg.Version}
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

// Folder defines filesystem folder descriptor (path and content structure)
type Folder struct {
	Content *FolderContent
	Path    string
}

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

// VisualStudioSolution defines VS solution that contains *solution.Solution
// and it's path
type VisualStudioSolution struct {
	// Solution structure
	Solution *solution.Solution

	// filesystem path
	Path string
}

// ProjectHandler defines project handler function prototype
type ProjectHandler func(*MsbuildProject, *Folder)

// StringDecorator defines string decorating function
type StringDecorator func(s string) string

// PassThrough just return original string without modification
// it can be used as non destructive decorator
func PassThrough(s string) string {
	return s
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

type packages struct {
	XMLName  xml.Name       `xml:"packages"`
	Packages []nugetPackage `xml:"package"`
}

type nugetPackage struct {
	ID                    string `xml:"id,attr"`
	Version               string `xml:"version,attr"`
	TargetFramework       string `xml:"targetFramework,attr"`
	DevelopmentDependency string `xml:"developmentDependency,attr"`
}

type msbuildProject struct {
	XMLName           xml.Name           `xml:"Project"`
	Sdk               string             `xml:"Sdk,attr"`
	ToolsVersion      string             `xml:"ToolsVersion,attr"`
	DefaultTargets    string             `xml:"DefaultTargets,attr"`
	ID                string             `xml:"PropertyGroup>ProjectGuid"`
	Compiles          []include          `xml:"ItemGroup>Compile"`
	CLCompiles        []include          `xml:"ItemGroup>ClCompile"`
	CLInclude         []include          `xml:"ItemGroup>ClInclude"`
	Contents          []include          `xml:"ItemGroup>Content"`
	Nones             []include          `xml:"ItemGroup>None"`
	References        []reference        `xml:"ItemGroup>Reference"`
	ProjectReferences []projectReference `xml:"ItemGroup>ProjectReference"`
	PackageReferences []packageReference `xml:"ItemGroup>PackageReference"`
	OutputPaths       []string           `xml:"PropertyGroup>OutputPath"`
	Imports           []msbuildImport    `xml:"Import"`
}

type include struct {
	Include string `xml:"Include,attr"`
}

func (i *include) path() string {
	return solution.ToValidPath(i.Include)
}

type reference struct {
	Assembly string `xml:"Include,attr"`
	HintPath string `xml:"HintPath"`
}

type projectReference struct {
	Include     string `xml:"Include,attr"`
	ProjectGUID string `xml:"Project"`
	Name        string `xml:"Name"`
}

// Path gets referenced project path
func (r *projectReference) Path() string {
	return solution.ToValidPath(r.Include)
}

type packageReference struct {
	ID      string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

type msbuildImport struct {
	Project string `xml:"Project,attr"`
	Sdk     string `xml:"Sdk,attr"`
}

type filesystem struct {
	fs afero.Fs
}

func newFs(fs afero.Fs) scan.Filesystem {
	return &filesystem{fs: fs}
}

func (f *filesystem) Open(path string) (scan.File, error) {
	return f.fs.Open(path)
}
