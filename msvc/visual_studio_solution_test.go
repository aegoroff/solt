package msvc

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_AllProjectPaths(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	memfs := afero.NewMemMapFs()
	dir := "a/"
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(Vs2017), 0644)

	foldersTree := ReadSolutionDir(dir, memfs)

	solutions := SelectSolutions(foldersTree)

	// Act
	paths := solutions[0].AllProjectPaths(func(s string) string { return s })

	// Assert
	ass.Len(paths, 8)
	ass.True(len(paths) < len(solutions[0].Solution.Projects))
}

const Vs2017 = `
Microsoft Visual Studio Solution File, Format Version 12.00
# Visual Studio 15
VisualStudioVersion = 15.0.26403.0
MinimumVisualStudioVersion = 10.0.40219.1
Project("{930C7802-8A8C-48F9-8165-68863BCCD9DD}") = "logviewer.install", "logviewer.install\logviewer.install.wixproj", "{27060CA7-FB29-42BC-BA66-7FC80D498354}"
	ProjectSection(ProjectDependencies) = postProject
		{405827CB-84E1-46F3-82C9-D889892645AC} = {405827CB-84E1-46F3-82C9-D889892645AC}
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D} = {CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}
	EndProjectSection
EndProject
Project("{930C7802-8A8C-48F9-8165-68863BCCD9DD}") = "logviewer.install.bootstrap", "logviewer.install.bootstrap\logviewer.install.bootstrap.wixproj", "{1C0ED62B-D506-4E72-BBC2-A50D3926466E}"
	ProjectSection(ProjectDependencies) = postProject
		{27060CA7-FB29-42BC-BA66-7FC80D498354} = {27060CA7-FB29-42BC-BA66-7FC80D498354}
	EndProjectSection
EndProject
Project("{2150E333-8FDC-42A3-9474-1A3956D46DE8}") = "solution items", "solution items", "{3B960F8F-AD5D-45E7-92C0-05B65E200AC4}"
	ProjectSection(SolutionItems) = preProject
		.editorconfig = .editorconfig
		appveyor.yml = appveyor.yml
		logviewer.xml = logviewer.xml
		WiX.msbuild = WiX.msbuild
	EndProjectSection
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "logviewer.tests", "logviewer.tests\logviewer.tests.csproj", "{939DD379-CDC8-47EF-8D37-0E5E71D99D30}"
	ProjectSection(ProjectDependencies) = postProject
		{383C08FC-9CAC-42E5-9B02-471561479A74} = {383C08FC-9CAC-42E5-9B02-471561479A74}
	EndProjectSection
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "logviewer.logic", "logviewer.logic\logviewer.logic.csproj", "{383C08FC-9CAC-42E5-9B02-471561479A74}"
EndProject
Project("{2150E333-8FDC-42A3-9474-1A3956D46DE8}") = ".nuget", ".nuget", "{B720ED85-58CF-4840-B1AE-55B0049212CC}"
	ProjectSection(SolutionItems) = preProject
		.nuget\NuGet.Config = .nuget\NuGet.Config
	EndProjectSection
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "logviewer.engine", "logviewer.engine\logviewer.engine.csproj", "{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}"
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "logviewer.install.mca", "logviewer.install.mca\logviewer.install.mca.csproj", "{405827CB-84E1-46F3-82C9-D889892645AC}"
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "logviewer.ui", "logviewer.ui\logviewer.ui.csproj", "{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}"
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "logviewer.bench", "logviewer.bench\logviewer.bench.csproj", "{75E0C034-44C8-461B-A677-9A19566FE393}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Debug|Mixed Platforms = Debug|Mixed Platforms
		Debug|x86 = Debug|x86
		Release|Any CPU = Release|Any CPU
		Release|Mixed Platforms = Release|Mixed Platforms
		Release|x86 = Release|x86
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|Any CPU.ActiveCfg = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|Any CPU.Build.0 = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|Mixed Platforms.ActiveCfg = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|Mixed Platforms.Build.0 = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|x86.ActiveCfg = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Debug|x86.Build.0 = Debug|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|Any CPU.ActiveCfg = Release|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|Any CPU.Build.0 = Release|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|Mixed Platforms.ActiveCfg = Release|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|Mixed Platforms.Build.0 = Release|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|x86.ActiveCfg = Release|x86
		{27060CA7-FB29-42BC-BA66-7FC80D498354}.Release|x86.Build.0 = Release|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|Any CPU.ActiveCfg = Debug|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|Any CPU.Build.0 = Debug|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|Mixed Platforms.ActiveCfg = Debug|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|Mixed Platforms.Build.0 = Debug|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|x86.ActiveCfg = Debug|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Debug|x86.Build.0 = Debug|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|Any CPU.ActiveCfg = Release|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|Any CPU.Build.0 = Release|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|Mixed Platforms.ActiveCfg = Release|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|Mixed Platforms.Build.0 = Release|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|x86.ActiveCfg = Release|x86
		{1C0ED62B-D506-4E72-BBC2-A50D3926466E}.Release|x86.Build.0 = Release|x86
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Debug|Mixed Platforms.ActiveCfg = Debug|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Debug|Mixed Platforms.Build.0 = Debug|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Debug|x86.ActiveCfg = Debug|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Release|Any CPU.Build.0 = Release|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Release|Mixed Platforms.ActiveCfg = Release|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Release|Mixed Platforms.Build.0 = Release|Any CPU
		{939DD379-CDC8-47EF-8D37-0E5E71D99D30}.Release|x86.ActiveCfg = Release|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Debug|Mixed Platforms.ActiveCfg = Debug|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Debug|Mixed Platforms.Build.0 = Debug|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Debug|x86.ActiveCfg = Debug|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Release|Any CPU.Build.0 = Release|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Release|Mixed Platforms.ActiveCfg = Release|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Release|Mixed Platforms.Build.0 = Release|Any CPU
		{383C08FC-9CAC-42E5-9B02-471561479A74}.Release|x86.ActiveCfg = Release|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Debug|Mixed Platforms.ActiveCfg = Debug|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Debug|Mixed Platforms.Build.0 = Debug|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Debug|x86.ActiveCfg = Debug|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Release|Any CPU.Build.0 = Release|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Release|Mixed Platforms.ActiveCfg = Release|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Release|Mixed Platforms.Build.0 = Release|Any CPU
		{90E3A68D-C96D-4764-A1D0-F73D9F474BE4}.Release|x86.ActiveCfg = Release|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Debug|Mixed Platforms.ActiveCfg = Debug|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Debug|Mixed Platforms.Build.0 = Debug|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Debug|x86.ActiveCfg = Debug|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Release|Any CPU.Build.0 = Release|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Release|Mixed Platforms.ActiveCfg = Release|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Release|Mixed Platforms.Build.0 = Release|Any CPU
		{405827CB-84E1-46F3-82C9-D889892645AC}.Release|x86.ActiveCfg = Release|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Debug|Mixed Platforms.ActiveCfg = Debug|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Debug|Mixed Platforms.Build.0 = Debug|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Debug|x86.ActiveCfg = Debug|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Release|Any CPU.Build.0 = Release|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Release|Mixed Platforms.ActiveCfg = Release|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Release|Mixed Platforms.Build.0 = Release|Any CPU
		{CFBAE2FB-6E3F-44CF-9FC9-372D6EA8DD3D}.Release|x86.ActiveCfg = Release|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Debug|Mixed Platforms.ActiveCfg = Debug|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Debug|Mixed Platforms.Build.0 = Debug|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Debug|x86.ActiveCfg = Debug|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Debug|x86.Build.0 = Debug|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Release|Any CPU.Build.0 = Release|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Release|Mixed Platforms.ActiveCfg = Release|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Release|Mixed Platforms.Build.0 = Release|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Release|x86.ActiveCfg = Release|Any CPU
		{75E0C034-44C8-461B-A677-9A19566FE393}.Release|x86.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
EndGlobal
`
