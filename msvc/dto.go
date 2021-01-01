package msvc

import (
	"encoding/xml"
	"path/filepath"
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

type includes []include

type msbuildProject struct {
	XMLName           xml.Name           `xml:"Project"`
	Sdk               string             `xml:"Sdk,attr"`
	ToolsVersion      string             `xml:"ToolsVersion,attr"`
	DefaultTargets    string             `xml:"DefaultTargets,attr"`
	ID                string             `xml:"PropertyGroup>ProjectGuid"`
	Compiles          includes           `xml:"ItemGroup>Compile"`
	CLCompiles        includes           `xml:"ItemGroup>ClCompile"`
	CLInclude         includes           `xml:"ItemGroup>ClInclude"`
	Contents          includes           `xml:"ItemGroup>Content"`
	Nones             includes           `xml:"ItemGroup>None"`
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

func (in *includes) paths(basePath string) []string {
	result := make([]string, len(*in))

	for i, c := range *in {
		result[i] = filepath.Join(basePath, c.path())
	}

	return result
}
