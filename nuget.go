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
