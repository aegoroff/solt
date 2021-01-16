package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/internal/out"
	"solt/internal/sys"
	"testing"
)

func Test_InfoCmd_InfoAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(wixSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a1.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"w/w.wixproj", []byte(wix), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "in", dir)

	// Assert
	actual := env.String()
	ass.Equal(sys.ToValidPath(` a/a.sln
  Header                           Microsoft Visual Studio Solution File, Format Version 12.00
  Product                          # Visual Studio Version 16
  Visual Studio Version            16.0.30104.148
  Minimum Visual Studio Version    10.0.40219.1

   Project type                             Count
   ------------                             -----
   WiX (Windows Installer XML)              1
   C#                                       2
   {00000000-000-0000-0000-000000000000}    1

   Configuration
   -------------
   Debug
   Release

   Platform
   --------
   Any CPU

 a/a1.sln
  Header                           Microsoft Visual Studio Solution File, Format Version 12.00
  Product                          # Visual Studio Version 16
  Visual Studio Version            16.0.30104.148
  Minimum Visual Studio Version    10.0.40219.1

   Project type    Count
   ------------    -----
   C#              1

   Configuration
   -------------
   Debug
   Release

   Platform
   --------
   Any CPU

 Totals:
  Solutions                                2
  Projects                                 5
                                           
  Project type                             Count    %         Solutions    %     
  ------------                             -----    ------    ---------    ------
  C#                                       3        60.00%    2            100.00%
  WiX (Windows Installer XML)              1        20.00%    1            50.00%
  {00000000-000-0000-0000-000000000000}    1        20.00%    1            50.00%
`), actual)
}

func Test_InfoNoPath_OutputHelp(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	env := out.NewMemoryEnvironment()

	// Act
	_ = Execute(memfs, env, "in")

	// Assert
	actual := env.String()
	ass.Contains(actual, "Get information about found solutions")
}

const wixSolutionContent = `

Microsoft Visual Studio Solution File, Format Version 12.00
# Visual Studio Version 16
VisualStudioVersion = 16.0.30104.148
MinimumVisualStudioVersion = 10.0.40219.1
Project("{930C7802-8A8C-48F9-8165-68863BCCD9DD}") = "w", "w\w.wixproj", "{27060CA7-FB29-42BC-BA66-7FC80D498354}"
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "a", "a\a.csproj", "{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}"
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "b", "b\b.csproj", "{1C0ED62B-D506-4E72-BBC2-A50D3926466E}"
EndProject
Project("{00000000-000-0000-0000-000000000000}") = "test.prop", "b\test.prop", "{97E04FD2-B904-46BE-930A-B270D264E83C}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|Any CPU.ActiveCfg = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|Any CPU.Build.0 = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|Any CPU.ActiveCfg = Release|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|Any CPU.Build.0 = Release|x86
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{3F69AE61-CC2E-40DB-B3CE-77ABD9BF327F}.Release|Any CPU.Build.0 = Release|Any CPU
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|Any CPU.Build.0 = Release|Any CPU
		{97E04FD2-B904-46BE-930A-B270D264E83C}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{97E04FD2-B904-46BE-930A-B270D264E83C}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{97E04FD2-B904-46BE-930A-B270D264E83C}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{97E04FD2-B904-46BE-930A-B270D264E83C}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
	GlobalSection(ExtensibilityGlobals) = postSolution
		SolutionGuid = {99ED9FBC-1A1E-43E7-88E0-9B83F88985CC}
	EndGlobalSection
EndGlobal

`

const wix = `<?xml version="1.0" encoding="utf-8"?>
<Project ToolsVersion="4.0" DefaultTargets="Build" xmlns="http://schemas.microsoft.com/developer/msbuild/2003">
  <PropertyGroup>
    <Configuration Condition=" '$(Configuration)' == '' ">Debug</Configuration>
    <Platform Condition=" '$(Platform)' == '' ">x86</Platform>
    <ProductVersion>3.6</ProductVersion>
    <ProjectGuid>27060ca7-fb29-42bc-ba66-7fc80d498354</ProjectGuid>
    <SchemaVersion>2.0</SchemaVersion>
    <OutputName>logviewer.install</OutputName>
    <OutputType>Package</OutputType>
    <WixTargetsPath Condition=" '$(WixTargetsPath)' == '' AND '$(MSBuildExtensionsPath32)' != '' ">$(MSBuildExtensionsPath32)\Microsoft\WiX\v3.x\Wix.targets</WixTargetsPath>
    <WixTargetsPath Condition=" '$(WixTargetsPath)' == '' ">$(MSBuildExtensionsPath)\Microsoft\WiX\v3.x\Wix.targets</WixTargetsPath>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Debug|x86' ">
    <OutputPath>bin\$(Configuration)\</OutputPath>
    <IntermediateOutputPath>obj\$(Configuration)\</IntermediateOutputPath>
    <DefineConstants>Debug;CONFIGURATION=Debug</DefineConstants>
  </PropertyGroup>
  <PropertyGroup Condition=" '$(Configuration)|$(Platform)' == 'Release|x86' ">
    <OutputPath>bin\$(Configuration)\</OutputPath>
    <IntermediateOutputPath>obj\$(Configuration)\</IntermediateOutputPath>
    <DefineConstants>CONFIGURATION=Release</DefineConstants>
  </PropertyGroup>
  <ItemGroup>
    <Compile Include="files.wxs" />
    <Compile Include="DirectoriesDefinition.wxs" />
    <Compile Include="Product.wxs" />
  </ItemGroup>
  <ItemGroup>
    <WixExtension Include="WixUtilExtension">
      <HintPath>$(WixExtDir)\WixUtilExtension.dll</HintPath>
      <Name>WixUtilExtension</Name>
    </WixExtension>
    <WixExtension Include="WixUIExtension">
      <HintPath>$(WixExtDir)\WixUIExtension.dll</HintPath>
      <Name>WixUIExtension</Name>
    </WixExtension>
    <WixExtension Include="WixNetFxExtension">
      <HintPath>$(WixExtDir)\WixNetFxExtension.dll</HintPath>
      <Name>WixNetFxExtension</Name>
    </WixExtension>
  </ItemGroup>
  <ItemGroup>
    <Content Include="Variables.wxi" />
  </ItemGroup>
  <Import Project="$(WixTargetsPath)" />
  <Import Project="..\WiX.msbuild" />
  <Target Name="BeforeBuild" DependsOnTargets="UpdateWix">
    <ItemGroup>
      <BinFile Include="bin\logviewer\**\*.exe" />
      <BinFile Include="bin\logviewer\**\*.dll" />
    </ItemGroup>
    <PropertyGroup Condition="$(KeyFile) != '' AND $(CertPassword) != '' AND $(SignTool) != '' AND Exists('$(KeyFile)')">
      <SignCommand>"$(SignTool)" sign /f "$(KeyFile)" /p $(CertPassword) /t http://timestamp.verisign.com/scripts/timestamp.dll /v /d "logviewer" /du https://github.com/aegoroff/logviewer %251</SignCommand>
    </PropertyGroup>
    <PropertyGroup Condition="$(SignCommand) != ''">
      <TmpFile>tmp.bat</TmpFile>
    </PropertyGroup>
    <WriteLinesToFile Condition="$(SignCommand) != ''" File="$(TmpFile)" Lines="$(SignCommand)" />
    <Exec Condition="$(SignCommand) != ''" Command="$(TmpFile) &quot;%(BinFile.Identity)&quot; &gt; NUL 2&gt;&amp;1" WorkingDirectory="$(MsBuildThisFileDirectory)" />
    <Delete Files="$(TmpFile)" Condition="$(SignCommand) != ''" />
  </Target>
  <Target Name="AfterBuild">
    <RemoveDir Directories="bin\logviewer" />
  </Target>
  <Target Name="CreateFiles">
    <ItemGroup>
      <PackageFile Include="$(MSBuildProjectDirectory)\..\logviewer.ui\bin\$(Configuration)\**\*.exe" Exclude="$(MSBuildProjectDirectory)\..\logviewer.ui\bin\$(Configuration)\**\*.vshost.exe" />
      <PackageFile Include="$(MSBuildProjectDirectory)\..\logviewer.ui\bin\$(Configuration)\**\*.dll" />
      <PackageFile Include="$(MSBuildProjectDirectory)\..\logviewer.ui\bin\$(Configuration)\**\*.pdb" />
      <PackageFile Include="$(MSBuildProjectDirectory)\..\logviewer.ui\bin\$(Configuration)\**\*.patterns" />
      <PackageFile Include="$(MSBuildProjectDirectory)\..\logviewer.ui\bin\$(Configuration)\**\*.config" Exclude="$(MSBuildProjectDirectory)\..\logviewer.ui\bin\$(Configuration)\**\*.vshost.exe.config" />
      <PackageFile Include="$(MSBuildProjectDirectory)\..\LICENSE.txt" />
    </ItemGroup>
    <MakeDir Directories="bin\logviewer" Condition="!Exists('bin\logviewer')" />
    <Copy SourceFiles="@(PackageFile)" DestinationFiles="@(PackageFile->'bin\logviewer\%(RecursiveDir)%(Filename)%(Extension)')" />
    <HeatDirectory ToolPath="$(WixInstallPath)" VerboseOutput="true" GenerateGuidsNow="true" OutputFile="$(ProjectDir)files.wxs" SuppressFragments="true" Directory="bin\logviewer" ComponentGroupName="logviewer" DirectoryRefId="INSTALLFOLDER" KeepEmptyDirectories="false" SuppressRootDirectory="true" PreprocessorVariable="var.SourcePath" SuppressRegistry="true" />
  </Target>
  <Target Name="UpdateWix" DependsOnTargets="CreateFiles">
    <ItemGroup>
      <RegexTransform Include="$(ProjectDir)files.wxs">
        <Find><![CDATA[File\s+Id="\w+"\s+KeyPath="yes"\s+Source="(.+)(logviewer\.ui\.exe|LICENSE\.txt)"]]></Find>
        <ReplaceWith><![CDATA[File Id="$2" KeyPath="yes" Source="$1$2"]]></ReplaceWith>
      </RegexTransform>
      <RegexTransform Include="$(ProjectDir)files.wxs">
        <Find><![CDATA[<Wix(.+?)>(\s*)<Fragment>]]></Find>
        <ReplaceWith><![CDATA[<Wix$1>$2<%3Finclude "Variables.wxi" %3F>$2<Fragment>]]></ReplaceWith>
        <Options>Singleline</Options>
      </RegexTransform>
    </ItemGroup>
    <RegexTransform Items="@(RegexTransform)" />
  </Target>
</Project>`
