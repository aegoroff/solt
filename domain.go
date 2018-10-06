package main

import "encoding/xml"

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
    Id                string             `xml:"PropertyGroup>ProjectGuid"`
    Compiles          []Include          `xml:"ItemGroup>Compile"`
    CLCompiles        []Include          `xml:"ItemGroup>ClCompile"`
    CLInclude         []Include          `xml:"ItemGroup>ClInclude"`
    Contents          []Include          `xml:"ItemGroup>Content"`
    Nones             []Include          `xml:"ItemGroup>None"`
    References        []Reference        `xml:"ItemGroup>Reference"`
    ProjectReferences []ProjectReference `xml:"ItemGroup>ProjectReference"`
    OutputPaths       []string           `xml:"PropertyGroup>OutputPath"`
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
