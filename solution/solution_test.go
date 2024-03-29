package solution

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_ParseSolution_ParsedSolution(t *testing.T) {
	var tests = []struct {
		input                string
		expectedComment      string
		expectedHead         string
		expectedVer          string
		expectedMinVer       string
		expectedProjectCount int
		expectedGsCount      int
		expectedProjectType  string
	}{
		{Vs2017, "# Visual Studio 15", "Microsoft Visual Studio Solution File, Format Version 12.00", "15.0.26403.0", "10.0.40219.1", 10, 3, "WiX (Windows Installer XML)"},
		{unix2Win(Vs2017), "# Visual Studio 15", "Microsoft Visual Studio Solution File, Format Version 12.00", "15.0.26403.0", "10.0.40219.1", 10, 3, "WiX (Windows Installer XML)"},
		{Vs2013, "# Visual Studio 2013", "Microsoft Visual Studio Solution File, Format Version 12.00", "12.0.31101.0", "10.0.40219.1", 1, 3, "C#"},
		{Vs2010, "# Visual Studio 2010", "Microsoft Visual Studio Solution File, Format Version 11.00", "", "", 1, 3, "C#"},
		{Vs2008, "# Visual Studio 2008", "Microsoft Visual Studio Solution File, Format Version 10.00", "", "", 1, 3, "C#"},
		{Vs2008StartsWithComment, "# Visual Studio 2008", "Microsoft Visual Studio Solution File, Format Version 10.00", "", "", 1, 3, "C#"},
		{Vs2008StartsWithCommentAndEmptyAfterIt, "# Visual Studio 2008", "Microsoft Visual Studio Solution File, Format Version 10.00", "", "", 1, 3, "C#"},
		{Vs7, "", "Microsoft Visual Studio Solution File, Format Version 7.00", "", "", 3, 6, "C#"},
	}

	// Act
	for _, test := range tests {
		t.Run(test.expectedComment, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)
			sol := parse(test.input, false)

			// Assert
			ass.Equal(test.expectedComment, sol.Comment)
			ass.Equal(test.expectedHead, sol.Header)
			ass.Equal(test.expectedVer, sol.VisualStudioVersion)
			ass.Equal(test.expectedMinVer, sol.MinimumVisualStudioVersion)
			ass.Equal(test.expectedProjectCount, len(sol.Projects))
			ass.Equal(test.expectedGsCount, len(sol.GlobalSections))
			ass.Equal(test.expectedProjectType, sol.Projects[0].Type)
		})
	}
}

func Test_ParseInvalidSolution_NoCrashHeadExtracted(t *testing.T) {
	// Arrange
	ass := assert.New(t)

	// Act
	sol := parse(Invalid, true)

	// Assert
	ass.Equal("# Visual Studio 2013", sol.Comment)
}

func Test_ParseInvalidSolution(t *testing.T) {
	var tests = []struct {
		name  string
		input string
	}{
		{"1", "\nMicrosoft Visual Studio Solution File, Format Version 12.00\n# Visual Studio 2013\nVisualStudioVersion = 12.0.31101.0\nMinimumVisualStudioVersion = 10.0.40219.1\nProject(\"{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}\") = \"Grok\", \"Grok\\Grok.csproj\", \"{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}\"\nEndProject\nGlobal\n\tGlobalSection(SolutionConfigurationPlatforms) = preSolution\n\t\tDebug|Any CPU = Debug|Any CPU\n\t\tRelease|Any CPU = Release|Any CPU\n\tEnnGlobalSectionease|Any CPU = Release|Any CPU\\n\\tEnnGlobalSection\\n\\tGlobalSection(ProjectConfigurationPlatforms) = postSolution\\n\\t\\t{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Debug|Any CPU.ActiveCfg = Debug|Ady CPU\\n\\t\\t{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Debug|Any CPU."},
		{"2", "\nMicrosoft Visual Studio Solution File, Format Version 12.00\n# Visual Studio 2013\nVisualStudioVersion = 12.0.31101.0\nMinimumVisualStudioVersion = 10.0.40219.1\nProject(\"{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}\") = \"Grok\", \"Grok\\Grok.csproj\", \"{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}\"\nEndProject\nGlobal\n\tGlobalSection(SolutionConfigurationPlatforms) = preSolution\n\t\tDebug|Any CPU = Deb"},
		{"3", "A\n\t = \n0"},
		{"4", "A00000000000000000000000000000000\nA00000000\n\t = \n"},
		{"5", "A(\"{\")=\"\",\"\",\"{\"\n\t = =0"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ass := assert.New(t)

			// Act
			sol := parse(test.input, false)

			// Assert
			ass.NotNil(sol)
		})
	}
}

func unix2Win(s string) string {
	return strings.ReplaceAll(s, "\n", "\r\n")
}

const Vs2013 = `
Microsoft Visual Studio Solution File, Format Version 12.00
# Visual Studio 2013
VisualStudioVersion = 12.0.31101.0
MinimumVisualStudioVersion = 10.0.40219.1
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "Grok", "Grok\Grok.csproj", "{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
EndGlobal
`

const Vs2010 = `
Microsoft Visual Studio Solution File, Format Version 11.00
# Visual Studio 2010
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "NlogWrapper", "NlogWrapper\NlogWrapper.csproj", "{3CFD766D-B991-4444-84A6-49BDFF81CF63}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{3CFD766D-B991-4444-84A6-49BDFF81CF63}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{3CFD766D-B991-4444-84A6-49BDFF81CF63}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{3CFD766D-B991-4444-84A6-49BDFF81CF63}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{3CFD766D-B991-4444-84A6-49BDFF81CF63}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
EndGlobal
`

const Vs2008 = `
Microsoft Visual Studio Solution File, Format Version 10.00
# Visual Studio 2008
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "DataVirtualization", "DataVirtualization\DataVirtualization.csproj", "{8102706C-AA37-4250-8889-1240FEB6F92F}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
EndGlobal
`

const Vs2008StartsWithComment = `
# comment
Microsoft Visual Studio Solution File, Format Version 10.00
# Visual Studio 2008
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "DataVirtualization", "DataVirtualization\DataVirtualization.csproj", "{8102706C-AA37-4250-8889-1240FEB6F92F}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
EndGlobal
`

const Vs2008StartsWithCommentAndEmptyAfterIt = `
# comment

Microsoft Visual Studio Solution File, Format Version 10.00
# Visual Studio 2008
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "DataVirtualization", "DataVirtualization\DataVirtualization.csproj", "{8102706C-AA37-4250-8889-1240FEB6F92F}"
EndProject
Global
	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{8102706C-AA37-4250-8889-1240FEB6F92F}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
EndGlobal
`

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

const Vs7 = `Microsoft Visual Studio Solution File, Format Version 7.00
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "log4net.Ext.Trace", "log4net.Ext.Trace\cs\src\log4net.Ext.Trace.csproj", "{8C73DF1C-AB2B-4309-A3EC-1ED594239E15}"
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "log4net.Ext.EventID", "log4net.Ext.EventID\cs\src\log4net.Ext.EventID.csproj", "{CB985027-C009-4C0F-88C1-8CF11912EE4C}"
EndProject
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}") = "log4net.Ext.MarshalByRef", "log4net.Ext.MarshalByRef\cs\src\log4net.Ext.MarshalByRef.csproj", "{CB985027-C009-4C0F-88C1-8CF11912AE4C}"
EndProject
Global
	GlobalSection(DPCodeReviewSolutionGUID) = preSolution
		DPCodeReviewSolutionGUID = {00000000-0000-0000-0000-000000000000}
	EndGlobalSection
	GlobalSection(SolutionConfiguration) = preSolution
		ConfigName.0 = Debug
		ConfigName.1 = Release
		ConfigName.2 = ReleaseStrong
	EndGlobalSection
	GlobalSection(ProjectDependencies) = postSolution
	EndGlobalSection
	GlobalSection(ProjectConfiguration) = postSolution
		{8C73DF1C-AB2B-4309-A3EC-1ED594239E15}.Debug.ActiveCfg = Debug|.NET
		{8C73DF1C-AB2B-4309-A3EC-1ED594239E15}.Debug.Build.0 = Debug|.NET
		{8C73DF1C-AB2B-4309-A3EC-1ED594239E15}.Release.ActiveCfg = Release|.NET
		{8C73DF1C-AB2B-4309-A3EC-1ED594239E15}.Release.Build.0 = Release|.NET
		{8C73DF1C-AB2B-4309-A3EC-1ED594239E15}.ReleaseStrong.ActiveCfg = ReleaseStrong|.NET
		{8C73DF1C-AB2B-4309-A3EC-1ED594239E15}.ReleaseStrong.Build.0 = ReleaseStrong|.NET
		{CB985027-C009-4C0F-88C1-8CF11912EE4C}.Debug.ActiveCfg = Debug|.NET
		{CB985027-C009-4C0F-88C1-8CF11912EE4C}.Debug.Build.0 = Debug|.NET
		{CB985027-C009-4C0F-88C1-8CF11912EE4C}.Release.ActiveCfg = Release|.NET
		{CB985027-C009-4C0F-88C1-8CF11912EE4C}.Release.Build.0 = Release|.NET
		{CB985027-C009-4C0F-88C1-8CF11912EE4C}.ReleaseStrong.ActiveCfg = ReleaseStrong|.NET
		{CB985027-C009-4C0F-88C1-8CF11912EE4C}.ReleaseStrong.Build.0 = ReleaseStrong|.NET
		{CB985027-C009-4C0F-88C1-8CF11912AE4C}.Debug.ActiveCfg = Debug|.NET
		{CB985027-C009-4C0F-88C1-8CF11912AE4C}.Debug.Build.0 = Debug|.NET
		{CB985027-C009-4C0F-88C1-8CF11912AE4C}.Release.ActiveCfg = Release|.NET
		{CB985027-C009-4C0F-88C1-8CF11912AE4C}.Release.Build.0 = Release|.NET
		{CB985027-C009-4C0F-88C1-8CF11912AE4C}.ReleaseStrong.ActiveCfg = Release|.NET
		{CB985027-C009-4C0F-88C1-8CF11912AE4C}.ReleaseStrong.Build.0 = Release|.NET
	EndGlobalSection
	GlobalSection(ExtensibilityGlobals) = postSolution
	EndGlobalSection
	GlobalSection(ExtensibilityAddIns) = postSolution
	EndGlobalSection
EndGlobal
`

const Invalid = `
Microsoft Visual Studio Solution File
# Visual Studio 2013
VisualStudioVersion = 12.0.31101.0
MinimumVisualStudioVersion = 10.0.40219.1
Project("{FAE04EC0-301F-11D3-BF4B-00C04F79EFBC}")  "Grok", "Grok\Grok.csproj", "{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}"
EndProject

	GlobalSection(SolutionConfigurationPlatforms) = preSolution
		Debug|Any CPU = Debug|Any CPU
		Release|Any CPU = Release|Any CPU
	EndGlobalSection
	GlobalSection(ProjectConfigurationPlatforms) = postSolution
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Debug|Any CPU.ActiveCfg = Debug|Any CPU
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Debug|Any CPU.Build.0 = Debug|Any CPU
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Release|Any CPU.ActiveCfg = Release|Any CPU
		{EC6D1E9B-2DA0-4225-9109-E9CF1C924116}.Release|Any CPU.Build.0 = Release|Any CPU
	EndGlobalSection
	GlobalSection(SolutionProperties) = preSolution
		HideSolutionNode = FALSE
	EndGlobalSection
EndGlobal
`
