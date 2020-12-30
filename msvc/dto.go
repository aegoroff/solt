package msvc

import (
	"encoding/xml"
	"solt/solution"
)

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

type reference struct {
	Assembly string `xml:"Include,attr"`
	HintPath string `xml:"HintPath"`
}

type projectReference struct {
	Include     string `xml:"Include,attr"`
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

func (i *include) path() string {
	return solution.ToValidPath(i.Include)
}

// Path gets referenced project path
func (r *projectReference) Path() string {
	return solution.ToValidPath(r.Include)
}

func (p *packages) nugetPackages() []*NugetPackage {
	result := make([]*NugetPackage, len(p.Packages))

	for i, pkg := range p.Packages {
		result[i] = &NugetPackage{ID: pkg.ID, Version: pkg.Version}
	}

	return result
}
