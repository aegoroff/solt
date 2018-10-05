package main

import "encoding/xml"

type Packages struct {
    XMLName  xml.Name  `xml:"packages"`
    Packages []Package `xml:"package"`
}

type Package struct {
    Id                    string `xml:"id,attr"`
    Version               string `xml:"version,attr"`
    TargetFramework       string `xml:"targetFramework,attr"`
    DevelopmentDependency string `xml:"developmentDependency,attr"`
}

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

type Include struct {
    Path string `xml:"Include,attr"`
}

type Reference struct {
    Assembly string `xml:"Include,attr"`
    HintPath string `xml:"HintPath"`
}

type ProjectReference struct {
    Path        string `xml:"Include,attr"`
    ProjectGuid string `xml:"Project"`
    Name        string `xml:"Name"`
}
