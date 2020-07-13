package msvc

import (
	"encoding/xml"
	"github.com/akutz/sortfold"
	"github.com/google/btree"
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

// Less tests whether the current item is less than the given argument.
//
// This must provide a strict weak ordering.
// If !a.Less(b) && !b.Less(a), we treat this to mean a == b (i.e. we can only
// hold one of either a or b in the tree).
func (x *Folder) Less(y btree.Item) bool {
	return sortfold.CompareFold(x.Path, (y.(*Folder)).Path) < 0
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
	Path string `xml:"Include,attr"`
}

type reference struct {
	Assembly string `xml:"Include,attr"`
	HintPath string `xml:"HintPath"`
}

type projectReference struct {
	Path        string `xml:"Include,attr"`
	ProjectGUID string `xml:"Project"`
	Name        string `xml:"Name"`
}

type packageReference struct {
	ID      string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

type msbuildImport struct {
	Project string `xml:"Project,attr"`
	Sdk     string `xml:"Sdk,attr"`
}
