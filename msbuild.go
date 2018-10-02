package main

import "encoding/xml"

type Project struct {
    XMLName           xml.Name           `xml:"Project"`
    Id                string             `xml:"PropertyGroup>ProjectGuid"`
    Compiles          []Compile          `xml:"ItemGroup>Compile"`
    Contents          []Content          `xml:"ItemGroup>Content"`
    Nones             []None             `xml:"ItemGroup>None"`
    References        []Reference        `xml:"ItemGroup>Reference"`
    ProjectReferences []ProjectReference `xml:"ItemGroup>ProjectReference"`
}

type Compile struct {
    Path string `xml:"Include,attr"`
}

type Content struct {
    Path string `xml:"Include,attr"`
}

type None struct {
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
