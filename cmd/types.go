package cmd

import (
	"encoding/xml"
)

// Packages is Nuget packages structure
type Packages struct {
	XMLName  xml.Name  `xml:"packages"`
	Packages []Package `xml:"package"`
}

// Package is Nuget package definition
type Package struct {
	Id                    string `xml:"id,attr"`
	Version               string `xml:"version,attr"`
	TargetFramework       string `xml:"targetFramework,attr"`
	DevelopmentDependency string `xml:"developmentDependency,attr"`
}

// Project is MSBuild project definition
type Project struct {
	XMLName           xml.Name           `xml:"Project"`
	Sdk               string             `xml:"Sdk,attr"`
	ToolsVersion      string             `xml:"ToolsVersion,attr"`
	DefaultTargets    string             `xml:"DefaultTargets,attr"`
	Id                string             `xml:"PropertyGroup>ProjectGuid"`
	Compiles          []Include          `xml:"ItemGroup>Compile"`
	CLCompiles        []Include          `xml:"ItemGroup>ClCompile"`
	CLInclude         []Include          `xml:"ItemGroup>ClInclude"`
	Contents          []Include          `xml:"ItemGroup>Content"`
	Nones             []Include          `xml:"ItemGroup>None"`
	References        []Reference        `xml:"ItemGroup>Reference"`
	ProjectReferences []ProjectReference `xml:"ItemGroup>ProjectReference"`
	PackageReferences []PackageReference `xml:"ItemGroup>PackageReference"`
	OutputPaths       []string           `xml:"PropertyGroup>OutputPath"`
	Imports           []Import           `xml:"Import"`
}

// Include attribute in MSBuild file
type Include struct {
	Path string `xml:"Include,attr"`
}

// Reference definition in MSBuild file
type Reference struct {
	Assembly string `xml:"Include,attr"`
	HintPath string `xml:"HintPath"`
}

// ProjectReference is project reference definition in MSBuild file
type ProjectReference struct {
	Path        string `xml:"Include,attr"`
	ProjectGuid string `xml:"Project"`
	Name        string `xml:"Name"`
}

// PackageReference is nuget reference definition in MSBuild file
type PackageReference struct {
	Id      string `xml:"Include,attr"`
	Version string `xml:"Version,attr"`
}

// Import attribute in MSBuild file
type Import struct {
	Project string `xml:"Project,attr"`
	Sdk     string `xml:"Sdk,attr"`
}

type nugetPackage struct {
	Id      string
	Version string
}

type walkEntry struct {
	Size   int64
	Parent string
	Name   string
	IsDir  bool
}

func (p *Project) isSdkProject() bool {
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
