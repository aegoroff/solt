package solution

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

func Test_ParseSolution_ParsedSolution(t *testing.T) {
    // Arrange
    ass := assert.New(t)
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
        {Vs2013, "# Visual Studio 2013", "Microsoft Visual Studio Solution File, Format Version 12.00", "12.0.31101.0", "10.0.40219.1", 1, 3, "C#"},
        {Vs2010, "# Visual Studio 2010", "Microsoft Visual Studio Solution File, Format Version 11.00", "", "", 1, 3, "C#"},
        {Vs2008, "# Visual Studio 2008", "Microsoft Visual Studio Solution File, Format Version 10.00", "", "", 1, 3, "C#"},
    }

    // Act
    for _, test := range tests {
        sol := parse(test.input)

        // Assert
        ass.Equal(test.expectedComment, sol.Comment)
        ass.Equal(test.expectedHead, sol.Header)
        ass.Equal(test.expectedVer, sol.VisualStudioVersion)
        ass.Equal(test.expectedMinVer, sol.MinimumVisualStudioVersion)
        ass.Equal(test.expectedProjectCount, len(sol.Projects))
        ass.Equal(test.expectedGsCount, len(sol.GlobalSections))
        ass.Equal(test.expectedProjectType, sol.Projects[0].Type)
    }
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
