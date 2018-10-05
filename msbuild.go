package main

import "encoding/xml"

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
