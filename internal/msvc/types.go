package msvc

import (
	"encoding/xml"
)

const (
	solutionFileExt    = ".sln"
	csharpProjectExt   = ".csproj"
	cppProjectExt      = ".vcxproj"
	packagesConfigFile = "packages.config"
)

// NugetPackage defines nuget package descriptor
type NugetPackage struct {
	ID      string
	Version string
}

// IsSdkProject gets whether a project is a the new VS 2017 or later project
func (p *msbuildProject) IsSdkProject() bool {
	if len(p.Sdk) > 0 {
		return true
	}
	if len(p.Imports) == 0 {
		return false
	}
	for _, imp := range p.Imports {
		if len(imp.Sdk) > 0 {
			return true
		}
	}
	return false
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
