package cmd

import (
	"bytes"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"solt/cmd/api"
	"solt/solution"
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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "p", "-p", dir)

	// Assert
	actual := w.String()
	ass.Equal(solution.ToValidPath(`
  <gray>a\a (packages.config)</>
  Package            Version
  -------            -------
  CmdLine            1.0.7.509
  Newtonsoft.Json    12.0.1
`), actual)
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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "p", "-p", dir)

	// Assert
	actual := w.String()
	ass.Equal(solution.ToValidPath(`
  <gray>a\a (a.csproj)</>
  Package              Version
  -------              -------
  CommandLineParser    2.8.0
`), actual)
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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "-p", dir, "-m")

	// Assert
	actual := w.String()
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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "-p", dir, "-m")

	// Assert
	actual := w.String()
	ass.Equal(solution.ToValidPath(` <red>Different nuget package's versions in the same solution found:</>

  <gray>a\a.sln</>
  Package              Version
  -------              -------
  CommandLineParser    2.7.0, 2.8.0
`), actual)
}

func Test_NugetCmdFindMismatchVerbose_OutputAsSpecified(t *testing.T) {
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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "-p", dir, "-m", "-v")

	// Assert
	actual := w.String()
	ass.Equal(solution.ToValidPath(` <red>Different nuget package's versions in the same solution found:</>

  <gray>a\a.sln</>
  Package              Version
  -------              -------
  CommandLineParser    2.7.0, 2.8.0

     <gray>Package: CommandLineParser</>
     Project    Version
     -------    -------
     a\a        2.8.0
     a\b        2.7.0
`), actual)
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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "-p", dir)

	// Assert
	actual := w.String()
	ass.Equal(solution.ToValidPath(`
  <gray>a\a.sln</>
  Package            Version
  -------            -------
  CmdLine            1.0.7.509
  Newtonsoft.Json    12.0.1
`), actual)
}

func Test_NugetCmdBySolutionNoPackages_NoOutput(t *testing.T) {
	// Arrange
	ass := assert.New(t)
	dir := "a/"
	memfs := afero.NewMemMapFs()
	_ = memfs.MkdirAll(dir+"a/Properties", 0755)
	_ = afero.WriteFile(memfs, dir+"a.sln", []byte(testSolutionContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/a.csproj", []byte(testProjectContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/App.config", []byte(appConfigContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Program.cs", []byte(codeFileContent), 0644)
	_ = afero.WriteFile(memfs, dir+"a/Properties/AssemblyInfo.cs", []byte(assemblyInfoContent), 0644)

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "-p", dir)

	// Assert
	actual := w.String()
	ass.Equal(``, actual)
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

	w := bytes.NewBufferString("")
	env := api.NewStringEnvironment(w)

	// Act
	_ = Execute(memfs, env, "nu", "-p", "d/")

	// Assert
	actual := w.String()
	ass.Contains(actual, solution.ToValidPath("d\\a1\\a.sln"))
	ass.Contains(actual, solution.ToValidPath("d\\a\\a.sln"))
	ass.Contains(actual, "CommandLineParser    2.8.0")
	ass.Contains(actual, "CmdLine            1.0.7.509")
	ass.Contains(actual, "Newtonsoft.Json    12.0.1")
}
