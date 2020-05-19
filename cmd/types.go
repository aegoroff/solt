package cmd

import (
	"encoding/xml"
	"github.com/aegoroff/godatastruct/rbtree"
	"github.com/akutz/sortfold"
	"strings"
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

type projectTreeNode struct {
	info *folderInfo
	path string
}

func (x projectTreeNode) LessThan(y interface{}) bool {
	return sortfold.CompareFold(x.path, (y.(projectTreeNode)).path) < 0
}

func (x projectTreeNode) EqualTo(y interface{}) bool {
	return strings.EqualFold(x.path, (y.(projectTreeNode)).path)
}

func newProjectTreeNode(path string, info *folderInfo) *rbtree.Comparable {
	var r rbtree.Comparable
	r = projectTreeNode{
		path: path,
		info: info,
	}
	return &r
}
