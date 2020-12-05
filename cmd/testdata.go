package cmd

const testSolutionContent = `

Microsoft Visual Studio Solution File, Format Version 12.00
# Visual Studio Version 16
VisualStudioVersion = 16.0.30104.148
MinimumVisualStudioVersion = 10.0.40219.1
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "a", "a\a.csproj", "{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
	GlobalSection(ExtensibilityGlobals) = postSolution
		SolutionGuid = {99ED9FBC-1A1E-43E7-88E0-9B83F88985CC}
	EndGlobalSection
EndGlobal

`
const testProjectContent = `<?xml version="1.0" encoding="utf-8"?>
<Project ToolsVersion="15.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <Import Project="$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props" Condition="Exists('$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props')" />
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectGuid>{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}</ProjectGuid>
    <OutputType>Exe</OutputType>
    <RootNamespace>a</RootNamespace>
    <AssemblyName>a</AssemblyName>
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
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <PlatformTarget>AnyCPU</PlatformTarget>
    <DebugType>pdbonly</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release\</OutputPath>
    <DefineConstants>TRACE</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
  </PropertyGroup>
  <ItemGroup>
    <Reference Include="CmdLine, Version=1.0.7.509, Culture=neutral, PublicKeyToken=e6ac9ead3d3b48e4, processorArchitecture=MSIL">
      <HintPath>..\packages\CmdLine.1.0.7.509\lib\net40-Client\CmdLine.dll</HintPath>
    </Reference>
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
    <None Include="packages.config" />
  </ItemGroup>
  <Import Project="$(MSBuildToolsPath)\Microsoft.CSharp.targets" />
</Project>
`

const testProjectContent2 = `<?xml version="1.0" encoding="utf-8"?>
<Project ToolsVersion="15.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <Import Project="$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props" Condition="Exists('$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props')" />
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectGuid>{7DBA74EF-3FCF-422B-B6EC-5EC4849C575A}</ProjectGuid>
    <OutputType>Exe</OutputType>
    <RootNamespace>braces</RootNamespace>
    <AssemblyName>braces</AssemblyName>
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
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <PlatformTarget>AnyCPU</PlatformTarget>
    <DebugType>pdbonly</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release\</OutputPath>
    <DefineConstants>TRACE</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
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

const testProjectContent3 = `<?xml version="1.0" encoding="utf-8"?>
<Project ToolsVersion="15.0" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <Import Project="$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props" Condition="Exists('$(MSBuildExtensionsPath)\$(MSBuildToolsVersion)\Microsoft.Common.props')" />
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">AnyCPU</Platform>
    <ProjectGuid>{7DBA74EF-3FCF-422B-B6EC-5EC4849C575A}</ProjectGuid>
    <OutputType>Exe</OutputType>
    <RootNamespace>braces</RootNamespace>
    <AssemblyName>braces</AssemblyName>
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
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|AnyCPU' ">
    <PlatformTarget>AnyCPU</PlatformTarget>
    <DebugType>pdbonly</DebugType>
    <Optimize>true</Optimize>
    <OutputPath>bin\Release\</OutputPath>
    <DefineConstants>TRACE</DefineConstants>
    <ErrorReport>prompt</ErrorReport>
    <WarningLevel>4</WarningLevel>
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
    <Compile Include="..\a\Program.cs">
		<Link>Program.cs</Link>
    </Compile>
    <Compile Include="..\a\Properties\AssemblyInfo.cs">
		<Link>Properties\AssemblyInfo.cs</Link>
    </Compile>
  </ItemGroup>
  <ItemGroup>
    <None Include="App.config" />
  </ItemGroup>
  <Import Project="$(MSBuildToolsPath)\Microsoft.CSharp.targets" />
</Project>
`

const appConfigContent = `<?xml version="1.0" encoding="utf-8" ?>
<configuration>
    <startup> 
        <supportedRuntime version="v4.0" sku=".NETFramework,Version=v4.7.2" />
    </startup>
</configuration>
`
const codeFileContent = `
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace a
{
    class Program
    {
        static void Main(string[] args)
        {
        }
    }
}
`

const assemblyInfoContent = `
using System.Reflection;
using System.Runtime.CompilerServices;
using System.Runtime.InteropServices;

// General Information about an assembly is controlled through the following
// set of attributes. Change these attribute values to modify the information
// associated with an assembly.
[assembly: AssemblyTitle("a")]
[assembly: AssemblyDescription("")]
[assembly: AssemblyConfiguration("")]
[assembly: AssemblyCompany("")]
[assembly: AssemblyProduct("a")]
[assembly: AssemblyCopyright("Copyright ©  2020")]
[assembly: AssemblyTrademark("")]
[assembly: AssemblyCulture("")]

// Setting ComVisible to false makes the types in this assembly not visible
// to COM components.  If you need to access a type in this assembly from
// COM, set the ComVisible attribute to true on that type.
[assembly: ComVisible(false)]

// The following GUID is for the ID of the typelib if this project is exposed to COM
[assembly: Guid("3f69ae61-cc2e-40db-b3ce-77abd9bf327f")]

// Version information for an assembly consists of the following four values:
//
//      Major Version
//      Minor Version
//      Build Number
//      Revision
//
// You can specify all the values or you can default the Build and Revision Numbers
// by using the '*' as shown below:
// [assembly: AssemblyVersion("1.0.*")]
[assembly: AssemblyVersion("1.0.0.0")]
[assembly: AssemblyFileVersion("1.0.0.0")]
`
const packagesConfingContent = `<?xml version="1.0" encoding="utf-8"?>
<packages>
  <package id="CmdLine" version="1.0.7.509" targetFramework="net472" />
  <package id="Newtonsoft.Json" version="12.0.1" targetFramework="net472" />
</packages>`

const coreSolutionContent = `
Microsoft Visual Studio Solution File, Format Version 12.00
# Visual Studio Version 16
VisualStudioVersion = 16.0.30717.126
MinimumVisualStudioVersion = 10.0.40219.1
Project("{9A19103F-16F7-4668-BE54-9A1E7A4F7556}") = "a", "a\a.csproj", "{836F0FC6-D6A9-4EA1-8818-35C9BD545249}"
EndProject
Project("{9A19103F-16F7-4668-BE54-9A1E7A4F7556}") = "b", "b\b.csproj", "{295171B4-8CD0-4929-A66A-1752395121BA}"
EndProject
Project("{9A19103F-16F7-4668-BE54-9A1E7A4F7556}") = "c", "c\c.csproj", "{9C13E797-65B5-4CDA-BA15-AE260AAA23C6}"
EndProject
Project("{2150E333-8FDC-42A3-9474-1A3956D46DE8}") = "f", "f", "{97E04FD2-B904-46BE-930A-B270D264E83C}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{836F0FC6-D6A9-4EA1-8818-35C9BD545249}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{836F0FC6-D6A9-4EA1-8818-35C9BD545249}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{836F0FC6-D6A9-4EA1-8818-35C9BD545249}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{836F0FC6-D6A9-4EA1-8818-35C9BD545249}.Release|Any CPU.Build.0 = Release|Any CPU
		{295171B4-8CD0-4929-A66A-1752395121BA}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{295171B4-8CD0-4929-A66A-1752395121BA}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{295171B4-8CD0-4929-A66A-1752395121BA}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{295171B4-8CD0-4929-A66A-1752395121BA}.Release|Any CPU.Build.0 = Release|Any CPU
		{9C13E797-65B5-4CDA-BA15-AE260AAA23C6}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{9C13E797-65B5-4CDA-BA15-AE260AAA23C6}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{9C13E797-65B5-4CDA-BA15-AE260AAA23C6}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{9C13E797-65B5-4CDA-BA15-AE260AAA23C6}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
	GlobalSection(NestedProjects) = preSolution
		{295171B4-8CD0-4929-A66A-1752395121BA} = {97E04FD2-B904-46BE-930A-B270D264E83C}
		{9C13E797-65B5-4CDA-BA15-AE260AAA23C6} = {97E04FD2-B904-46BE-930A-B270D264E83C}
	EndGlobalSection
	GlobalSection(ExtensibilityGlobals) = postSolution
		SolutionGuid = {31AA0C31-1842-4E8E-ADE9-FFE92FD1C22F}
	EndGlobalSection
EndGlobal
`

const aSdkProjectContent = `<Project Sdk="Microsoft.NET.Sdk">

  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>netcoreapp3.1</TargetFramework>
  </PropertyGroup>

  <ItemGroup>
    <ProjectReference Include="..\b\b.csproj" />
    <ProjectReference Include="..\c\c.csproj" />
  </ItemGroup>

</Project>
`

const bSdkProjectContent = `<Project Sdk="Microsoft.NET.Sdk">

  <PropertyGroup>
    <TargetFramework>netstandard2.0</TargetFramework>
  </PropertyGroup>

</Project>
`

const cSdkProjectContent = `<Project Sdk="Microsoft.NET.Sdk">

  <PropertyGroup>
    <TargetFramework>netstandard2.0</TargetFramework>
  </PropertyGroup>

  <ItemGroup>
    <ProjectReference Include="..\b\b.csproj" />
  </ItemGroup>

</Project>
`
