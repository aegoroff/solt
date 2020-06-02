package cmd

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"solt/internal/sys"
	"strings"
	"testing"
)

func Test_sortAndOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items := []string{"rr", "aa", "xy"}
	buf := bytes.NewBufferString("")

	// Act
	sortAndOutput(buf, items)

	// Assert
	ass.Equal(" aa\n rr\n xy\n", buf.String())
}

func Test_outputSortedMap(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	items1 := []string{"rr", "gt", "xy"}
	items2 := []string{"ff", "lz", "xy"}

	m := map[string][]string{"a": items1, "b": items2}
	buf := bytes.NewBufferString("")

	// Act
	outputSortedMap(buf, m, "SI")

	// Assert
	ass.Equal("\nSI: a\n gt\n rr\n xy\n\nSI: b\n ff\n lz\n xy\n", buf.String())
}

func Test_unmarshalXmlPackagesConfig(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	packages := Packages{}
	const packangesconfig = `
<?xml version="1.0" encoding="utf-8"?>
<packages>
  <package id="YaccLexTools" version="0.2.2" targetFramework="net45" />
</packages>
`
	r := strings.NewReader(packangesconfig)
	// Act
	sys.UnmarshalXml(r, &packages)

	// Assert
	ass.Equal(1, len(packages.Packages))
	ass.Equal("YaccLexTools", packages.Packages[0].Id)
	ass.Equal("0.2.2", packages.Packages[0].Version)
	ass.Equal("net45", packages.Packages[0].TargetFramework)
}

func Test_unmarshalXmlMsbuildProject(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	prj := Project{}
	const project = `
<?xml version="1.0" encoding="utf-8"?>
<Project ToolsVersion="15.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <Import Project="$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props" Condition="Exists('$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props')" />
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectGuid>{99B7AE2B-EF73-48A6-BBE1-ACF5E0CA569D}</ProjectGuid>
    <OutputType>Exe</OutputType>
    <RootNamespace>sort</RootNamespace>
    <AssemblyName>sort</AssemblyName>
    <TargetFrameworkVersion>v4.7.2</TargetFrameworkVersion>
    <FileAlignment>512</FileAlignment>
    <AutoGenerateBindingRedirects>true</AutoGenerateBindingRedirects>
    <Deterministic>true</Deterministic>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|AnyCPU' ">
    <PlatformTarget>AnyCPU</PlatformTarget>
    <DebugSymbols>true</DebugSymbols>
    <DebugType>full</DebugType>
    <Optimize>false</Optimize>
    <OutputPath>bin\Debug\</OutputPath>
    <DefineConstants>DEBUG;TRACE</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <AllowUnsafeBlocks>true</AllowUnsafeBlocks>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <PlatformTarget>AnyCPU</PlatformTarget>
    <DebugType>pdbonly</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release\</OutputPath>
    <DefineConstants>TRACE</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
    <AllowUnsafeBlocks>true</AllowUnsafeBlocks>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="System" />
    <Reference Include="System.Core" />
    <Reference Include="System.Xml.Linq" />
    <Reference Include="System.Data.DataSetExtensions" />
    <Reference Include="Microsoft.CSharp" />
    <Reference Include="System.Data" />
    <Reference Include="System.Net.Http" />
    <Reference Include="System.Xml" />
  </ItemGroup>
  <ItemGroup>
    <Compile Include="Program.cs" />
    <Compile Include="Properties\AssemblyInfo.cs" />
  </ItemGroup>
  <ItemGroup>
    <None Include="App.config" />
  </ItemGroup>
  <Import Project="$(MSBuildToolsPath)\Microsoft.CSharp.targets" />
</Project>
`
	r := strings.NewReader(project)

	// Act
	sys.UnmarshalXml(r, &prj)

	// Assert
	ass.Equal(2, len(prj.Compiles))
	ass.Equal(1, len(prj.Nones))
	ass.Equal("{99B7AE2B-EF73-48A6-BBE1-ACF5E0CA569D}", prj.Id)
	ass.Equal(8, len(prj.References))
}
