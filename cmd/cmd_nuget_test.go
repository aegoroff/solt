package cmd

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NugetCmd_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	_ = Execute("nu", "-p", dir, "-r")

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal(`
 a\a
  Package            Version
  -------            -------
  CmdLine            1.0.7.509
  Newtonsoft.Json    12.0.1
`, actual)
}

func Test_NugetCmdOnSdkProjects_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	_ = Execute("nu", "-p", dir, "-r")

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal(`
 a\a
  Package              Version
  -------              -------
  CommandLineParser    2.8.0
`, actual)
}

func Test_NugetCmdFindMismatchNoMismath_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	_ = Execute("nu", "-p", dir, "-m")

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal("", actual)
}

func Test_NugetCmdFindMismatch_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectWithNugetContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	_ = Execute("nu", "-p", dir, "-m")

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal(` <red>Different nuget package's versions in the same solution found:</>
 a\a.sln
  Package              Version
  -------              -------
  CommandLineParser    2.7.0, 2.8.0
`, actual)
}

func Test_NugetCmdBySolution_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	_ = Execute("nu", "-p", dir)

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal(`
 a\a.sln
  Package            Version
  -------            -------
  CmdLine            1.0.7.509
  Newtonsoft.Json    12.0.1
`, actual)
}

func Test_NugetCmdBySolutionManySolutions_OutputAsSpecified(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "d/a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/packages.config", []byte(packagesConfingContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	dir = "d/a1/"
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(coreSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(aSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/b.csproj", []byte(bSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"b/Class1.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/c.csproj", []byte(cSdkProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"c/Class1.cs", []byte(codeFileContent), 0644)

	appPrinter = newMockPrn()
	appFileSystem = memfs

	// Act
	_ = Execute("nu", "-p", "d/")

	// Assert
	actual := appPrinter.(*mockprn).String()
	ass.Equal(`
 d\a1\a.sln
  Package              Version
  -------              -------
  CommandLineParser    2.8.0

 d\a\a.sln
  Package            Version
  -------            -------
  CmdLine            1.0.7.509
  Newtonsoft.Json    12.0.1
`, actual)
}
